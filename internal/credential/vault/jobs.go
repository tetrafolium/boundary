package vault

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/boundary/internal/db"
	"github.com/hashicorp/boundary/internal/errors"
	"github.com/hashicorp/boundary/internal/kms"
	"github.com/hashicorp/boundary/internal/scheduler"
	"github.com/hashicorp/go-hclog"
	vault "github.com/hashicorp/vault/api"
	ua "go.uber.org/atomic"
)

const (
	tokenRenewalJobName         = "vault_token_renewal"
	tokenRevocationJobName      = "vault_token_revocation"
	credentialRenewalJobName    = "vault_credential_renewal"
	credentialRevocationJobName = "vault_credential_revocation"

	defaultNextRunIn = 5 * time.Minute
	renewalWindow    = 10 * time.Minute
)

// TokenRenewalJob is the recurring job that renews credential store Vault tokens that
// are in the `current` and `maintaining` state.  The TokenRenewalJob is not thread safe,
// an attempt to Run the job concurrently will result in an JobAlreadyRunning error.
type TokenRenewalJob struct {
	reader db.Reader
	writer db.Writer
	kms    *kms.Kms
	logger hclog.Logger
	limit  int

	running      ua.Bool
	numTokens    int
	numProcessed int
}

// NewTokenRenewalJob creates a new in-memory TokenRenewalJob.
//
// WithLimit is the only supported option.
func NewTokenRenewalJob(r db.Reader, w db.Writer, kms *kms.Kms, logger hclog.Logger, opt ...Option) (*TokenRenewalJob, error) {
	const op = "vault.NewTokenRenewalJob"
	switch {
	case r == nil:
		return nil, errors.New(errors.InvalidParameter, op, "missing db.Reader")
	case w == nil:
		return nil, errors.New(errors.InvalidParameter, op, "missing db.Writer")
	case kms == nil:
		return nil, errors.New(errors.InvalidParameter, op, "missing kms")
	case logger == nil:
		return nil, errors.New(errors.InvalidParameter, op, "missing logger")
	}

	opts := getOpts(opt...)
	if opts.withLimit == 0 {
		// zero signals the boundary defaults should be used.
		opts.withLimit = db.DefaultLimit
	}
	return &TokenRenewalJob{
		reader: r,
		writer: w,
		kms:    kms,
		logger: logger,
		limit:  opts.withLimit,
	}, nil
}

// Status returns the current status of the token renewal job.  Total is the total number
// of tokens that are set to be renewed. Completed is the number of tokens already renewed.
func (r *TokenRenewalJob) Status() scheduler.JobStatus {
	return scheduler.JobStatus{
		Completed: r.numProcessed,
		Total:     r.numTokens,
	}
}

// Run queries the vault credential repo for tokens that need to be renewed, it then creates
// a vault client and renews each token.  Can not be run in parallel, if Run is invoked while
// already running an error with code JobAlreadyRunning will be returned.
func (r *TokenRenewalJob) Run(ctx context.Context) error {
	const op = "vault.(TokenRenewalJob).Run"
	if !r.running.CAS(r.running.Load(), true) {
		return errors.New(errors.JobAlreadyRunning, op, "job already running")
	}
	defer r.running.Store(false)

	// Verify context is not done before running
	if err := ctx.Err(); err != nil {
		return errors.Wrap(err, op)
	}

	var ps []*privateStore
	// Fetch all tokens that will reach their renewal point within the renewalWindow.
	// This is done to avoid constantly scheduling the token renewal job when there are multiple tokens
	// set to renew in sequence.
	err := r.reader.SearchWhere(ctx, &ps, `token_renewal_time < wt_add_seconds_to_now(?)`, []interface{}{renewalWindow.Seconds()}, db.WithLimit(r.limit))
	if err != nil {
		return errors.Wrap(err, op)
	}

	// Set numProcessed and numTokens for status report
	r.numProcessed, r.numTokens = 0, len(ps)

	for _, s := range ps {
		// Verify context is not done before renewing next token
		if err := ctx.Err(); err != nil {
			return errors.Wrap(err, op)
		}
		if err := r.renewToken(ctx, s); err != nil {
			r.logger.Error("error renewing token", "credential store id", s.StoreId, "token status", s.TokenStatus, "error", err)
		}
		r.numProcessed++
	}

	return nil
}

func (r *TokenRenewalJob) renewToken(ctx context.Context, s *privateStore) error {
	const op = "vault.(TokenRenewalJob).renewToken"
	databaseWrapper, err := r.kms.GetWrapper(ctx, s.ScopeId, kms.KeyPurposeDatabase)
	if err != nil {
		return errors.Wrap(err, op, errors.WithMsg("unable to get database wrapper"))
	}
	if err = s.decrypt(ctx, databaseWrapper); err != nil {
		return errors.Wrap(err, op)
	}

	token := s.token()
	if token == nil {
		// Store has no token to renew
		return nil
	}

	vc, err := s.client()
	if err != nil {
		return errors.Wrap(err, op)
	}

	var respErr *vault.ResponseError
	renewedToken, err := vc.renewToken()
	if ok := errors.As(err, &respErr); ok && respErr.StatusCode == http.StatusForbidden {
		// Vault returned a 403 when attempting a renew self, the token is either expired
		// or malformed.  Set status to "expired" so credentials created with token can be
		// cleaned up.
		query, values := token.updateStatusQuery(ExpiredToken)
		numRows, err := r.writer.Exec(ctx, query, values)
		if err != nil {
			return errors.Wrap(err, op)
		}
		if numRows != 1 {
			return errors.New(errors.Unknown, op, "token expired but failed to update repo")
		}
		if s.TokenStatus == string(CurrentToken) {
			r.logger.Info("Vault credential store current token has expired", "credential store id", s.StoreId)
		}
		return nil
	}
	if err != nil {
		return errors.Wrap(err, op, errors.WithMsg("unable to renew vault token"))
	}

	tokenExpires, err := renewedToken.TokenTTL()
	if err != nil {
		return errors.Wrap(err, op, errors.WithMsg("unable to get vault token expiration"))
	}

	token.expiration = tokenExpires
	query, values := token.updateExpirationQuery()
	numRows, err := r.writer.Exec(ctx, query, values)
	if err != nil {
		return errors.Wrap(err, op)
	}
	if numRows != 1 {
		return errors.New(errors.Unknown, op, "token renewed but failed to update repo")
	}

	return nil
}

// NextRunIn queries the vault credential repo to determine when the next token renewal job should run.
func (r *TokenRenewalJob) NextRunIn() (time.Duration, error) {
	const op = "vault.(TokenRenewalJob).NextRunIn"
	next, err := nextRenewal(r)
	if err != nil {
		return defaultNextRunIn, errors.Wrap(err, op)
	}

	return next, nil
}

func nextRenewal(j scheduler.Job) (time.Duration, error) {
	const op = "vault.nextRenewal"
	var query string
	var r db.Reader
	switch job := j.(type) {
	case *TokenRenewalJob:
		query = tokenRenewalNextRunInQuery
		r = job.reader
	case *CredentialRenewalJob:
		query = credentialRenewalNextRunInQuery
		r = job.reader
	default:
		return 0, errors.New(errors.Unknown, op, "unknown job")
	}

	rows, err := r.Query(context.Background(), query, nil)
	if err != nil {
		return 0, errors.Wrap(err, op)
	}
	defer rows.Close()

	for rows.Next() {
		type NextRenewal struct {
			RenewalIn time.Duration
		}
		var n NextRenewal
		err = r.ScanRows(rows, &n)
		if err != nil {
			return 0, errors.Wrap(err, op)
		}
		if n.RenewalIn < 0 {
			// If we are past the next renewal time, return 0 to schedule immediately
			return 0, nil
		}
		return n.RenewalIn * time.Second, nil
	}

	return defaultNextRunIn, nil
}

// Name is the unique name of the job.
func (r *TokenRenewalJob) Name() string {
	return tokenRenewalJobName
}

// Description is the human readable description of the job.
func (r *TokenRenewalJob) Description() string {
	return "Periodically renews Vault credential store tokens that are in a maintaining or current state."
}

// TokenRevocationJob is the recurring job that revokes credential store Vault tokens that
// are in the `maintaining` state and have no credentials being used by an active or pending session.
// The TokenRevocationJob is not thread safe, an attempt to Run the job concurrently will result in
// an JobAlreadyRunning error.
type TokenRevocationJob struct {
	reader db.Reader
	writer db.Writer
	kms    *kms.Kms
	logger hclog.Logger
	limit  int

	running      ua.Bool
	numTokens    int
	numProcessed int
}

// NewTokenRevocationJob creates a new in-memory TokenRevocationJob.
//
// WithLimit is the only supported option.
func NewTokenRevocationJob(r db.Reader, w db.Writer, kms *kms.Kms, logger hclog.Logger, opt ...Option) (*TokenRevocationJob, error) {
	const op = "vault.NewTokenRevocationJob"
	switch {
	case r == nil:
		return nil, errors.New(errors.InvalidParameter, op, "missing db.Reader")
	case w == nil:
		return nil, errors.New(errors.InvalidParameter, op, "missing db.Writer")
	case kms == nil:
		return nil, errors.New(errors.InvalidParameter, op, "missing kms")
	case logger == nil:
		return nil, errors.New(errors.InvalidParameter, op, "missing logger")
	}

	opts := getOpts(opt...)
	if opts.withLimit == 0 {
		// zero signals the boundary defaults should be used.
		opts.withLimit = db.DefaultLimit
	}
	return &TokenRevocationJob{
		reader: r,
		writer: w,
		kms:    kms,
		logger: logger,
		limit:  opts.withLimit,
	}, nil
}

// Status returns the current status of the token revocation job.  Total is the total number
// of tokens that are set to be revoked. Completed is the number of tokens already revoked.
func (r *TokenRevocationJob) Status() scheduler.JobStatus {
	return scheduler.JobStatus{
		Completed: r.numProcessed,
		Total:     r.numTokens,
	}
}

// Run queries the vault credential repo for tokens that need to be revoked, it then creates
// a vault client and revokes each token.  Can not be run in parallel, if Run is invoked while
// already running an error with code JobAlreadyRunning will be returned.
func (r *TokenRevocationJob) Run(ctx context.Context) error {
	const op = "vault.(TokenRevocationJob).Run"
	if !r.running.CAS(r.running.Load(), true) {
		return errors.New(errors.JobAlreadyRunning, op, "job already running")
	}
	defer r.running.Store(false)

	// Verify context is not done before running
	if err := ctx.Err(); err != nil {
		return errors.Wrap(err, op)
	}

	// Fetch all tokens in the maintaining state that have no credentials in an active state
	where := `
token_status = ?
  and token_hmac not in (
    select token_hmac from credential_vault_credential 
     where status = ?
)
`
	whereArgs := []interface{}{MaintainingToken, ActiveCredential}

	var ps []*privateStore
	err := r.reader.SearchWhere(ctx, &ps, where, whereArgs, db.WithLimit(r.limit))
	if err != nil {
		return errors.Wrap(err, op)
	}

	// Set numProcessed and numTokens for s report
	r.numProcessed, r.numTokens = 0, len(ps)
	for _, s := range ps {
		// Verify context is not done before renewing next token
		if err := ctx.Err(); err != nil {
			return errors.Wrap(err, op)
		}
		if err := r.revokeToken(ctx, s); err != nil {
			r.logger.Error("error revoking token", "credential store id", s.StoreId, "error", err)
		}
		r.numProcessed++
	}

	return nil
}

func (r *TokenRevocationJob) revokeToken(ctx context.Context, s *privateStore) error {
	const op = "vault.(TokenRevocationJob).revokeToken"
	databaseWrapper, err := r.kms.GetWrapper(ctx, s.ScopeId, kms.KeyPurposeDatabase)
	if err != nil {
		return errors.Wrap(err, op, errors.WithMsg("unable to get database wrapper"))
	}
	if err = s.decrypt(ctx, databaseWrapper); err != nil {
		return errors.Wrap(err, op)
	}

	token := s.token()
	if token == nil {
		// Store has no token to revoke
		return nil
	}

	vc, err := s.client()
	if err != nil {
		return errors.Wrap(err, op)
	}

	var respErr *vault.ResponseError
	err = vc.revokeToken()
	if ok := errors.As(err, &respErr); ok && respErr.StatusCode == http.StatusForbidden {
		// Vault returned a 403 when attempting a revoke self, the token is already expired.
		// Clobber error and set status to "revoked" below.
		err = nil
	}
	if err != nil {
		return errors.Wrap(err, op, errors.WithMsg("unable to revoke vault token"))
	}

	query, values := token.updateStatusQuery(RevokedToken)
	numRows, err := r.writer.Exec(ctx, query, values)
	if err != nil {
		return errors.Wrap(err, op)
	}
	if numRows != 1 {
		return errors.New(errors.Unknown, op, "token revoked but failed to update repo")
	}

	return nil
}

// NextRunIn determines when the next token revocation job should run.
func (r *TokenRevocationJob) NextRunIn() (time.Duration, error) {
	return defaultNextRunIn, nil
}

// Name is the unique name of the job.
func (r *TokenRevocationJob) Name() string {
	return tokenRevocationJobName
}

// Description is the human readable description of the job.
func (r *TokenRevocationJob) Description() string {
	return "Periodically revokes Vault credential store tokens that are in a maintaining state and have no active credentials associated."
}

// CredentialRenewalJob is the recurring job that renews Vault credentials issued to a session.
// The CredentialRenewalJob is not thread safe, an attempt to Run the job concurrently will result
// in an JobAlreadyRunning error.
type CredentialRenewalJob struct {
	reader db.Reader
	writer db.Writer
	kms    *kms.Kms
	logger hclog.Logger
	limit  int

	running      ua.Bool
	numCreds     int
	numProcessed int
}

// NewCredentialRenewalJob creates a new in-memory CredentialRenewalJob.
//
// WithLimit is the only supported option.
func NewCredentialRenewalJob(r db.Reader, w db.Writer, kms *kms.Kms, logger hclog.Logger, opt ...Option) (*CredentialRenewalJob, error) {
	const op = "vault.NewCredentialRenewalJob"
	switch {
	case r == nil:
		return nil, errors.New(errors.InvalidParameter, op, "missing db.Reader")
	case w == nil:
		return nil, errors.New(errors.InvalidParameter, op, "missing db.Writer")
	case kms == nil:
		return nil, errors.New(errors.InvalidParameter, op, "missing kms")
	case logger == nil:
		return nil, errors.New(errors.InvalidParameter, op, "missing logger")
	}

	opts := getOpts(opt...)
	if opts.withLimit == 0 {
		// zero signals the boundary defaults should be used.
		opts.withLimit = db.DefaultLimit
	}
	return &CredentialRenewalJob{
		reader: r,
		writer: w,
		kms:    kms,
		logger: logger,
		limit:  opts.withLimit,
	}, nil
}

// Status returns the current status of the credential renewal job.  Total is the total number
// of credentials that are set to be renewed.  Completed is the number of credential already renewed.
func (r *CredentialRenewalJob) Status() scheduler.JobStatus {
	return scheduler.JobStatus{
		Completed: r.numProcessed,
		Total:     r.numCreds,
	}
}

// Run queries the vault credential repo for credentials that need to be renewed, it then creates
// a vault client and renews each credential.  Can not be run in parallel, if Run is invoked while
// already running an error with code JobAlreadyRunning will be returned.
func (r *CredentialRenewalJob) Run(ctx context.Context) error {
	const op = "vault.(CredentialRenewalJob).Run"
	if !r.running.CAS(r.running.Load(), true) {
		return errors.New(errors.JobAlreadyRunning, op, "job already running")
	}
	defer r.running.Store(false)

	// Verify context is not done before running
	if err := ctx.Err(); err != nil {
		return errors.Wrap(err, op)
	}

	var creds []*privateCredential
	// Fetch all active credentials that will reach their renewal point within the renewalWindow.
	// This is done to avoid constantly scheduling the credential renewal job when there are
	// multiple credentials set to renew in sequence.
	err := r.reader.SearchWhere(ctx, &creds, `renewal_time < wt_add_seconds_to_now(?) and status = ?`, []interface{}{renewalWindow.Seconds(), ActiveCredential}, db.WithLimit(r.limit))
	if err != nil {
		return errors.Wrap(err, op)
	}

	// Set numProcessed and numTokens for status report
	r.numProcessed, r.numCreds = 0, len(creds)
	for _, c := range creds {
		// Verify context is not done before renewing next token
		if err := ctx.Err(); err != nil {
			return errors.Wrap(err, op)
		}

		if err := r.renewCred(ctx, c); err != nil {
			r.logger.Error("error renewing credential", "credential id", c.PublicId, "error", err)
		}

		r.numProcessed++
	}

	return nil
}

func (r *CredentialRenewalJob) renewCred(ctx context.Context, c *privateCredential) error {
	const op = "vault.(CredentialRenewalJob).renewCred"
	databaseWrapper, err := r.kms.GetWrapper(ctx, c.ScopeId, kms.KeyPurposeDatabase)
	if err != nil {
		return errors.Wrap(err, op, errors.WithMsg("unable to get database wrapper"))
	}
	if err = c.decrypt(ctx, databaseWrapper); err != nil {
		return errors.Wrap(err, op)
	}

	vc, err := c.client()
	if err != nil {
		return errors.Wrap(err, op)
	}
	cred := c.toCredential()

	var respErr *vault.ResponseError
	// Subtract last renewal time from previous expiration time to get lease duration
	leaseDuration := c.ExpirationTime.AsTime().Sub(c.LastRenewalTime.AsTime())
	renewedCred, err := vc.renewLease(c.ExternalId, leaseDuration)
	if ok := errors.As(err, &respErr); ok && respErr.StatusCode == http.StatusBadRequest {
		// Vault returned a 400 when attempting a renew lease, the lease is either expired
		// or the leaseId is malformed.  Set status to "expired".
		query, values := cred.updateStatusQuery(ExpiredCredential)
		numRows, err := r.writer.Exec(ctx, query, values)
		if err != nil {
			return errors.Wrap(err, op)
		}
		if numRows != 1 {
			return errors.New(errors.Unknown, op, "credential expired but failed to update repo")
		}
		return nil
	}
	if err != nil {
		return errors.Wrap(err, op, errors.WithMsg("unable to renew credential"))
	}

	cred.expiration = time.Duration(renewedCred.LeaseDuration) * time.Second
	query, values := cred.updateExpirationQuery()
	numRows, err := r.writer.Exec(ctx, query, values)
	if err != nil {
		return errors.Wrap(err, op)
	}
	if numRows != 1 {
		return errors.New(errors.Unknown, op, "credential renewed but failed to update repo")
	}

	return nil
}

// NextRunIn queries the vault credential repo to determine when the next credential renewal job should run.
func (r *CredentialRenewalJob) NextRunIn() (time.Duration, error) {
	const op = "vault.(CredentialRenewalJob).NextRunIn"
	next, err := nextRenewal(r)
	if err != nil {
		return defaultNextRunIn, errors.Wrap(err, op)
	}

	return next, nil
}

// Name is the unique name of the job.
func (r *CredentialRenewalJob) Name() string {
	return credentialRenewalJobName
}

// Description is the human readable description of the job.
func (r *CredentialRenewalJob) Description() string {
	return "Periodically renews Vault credentials that are attached to an active/pending session (in the active state)."
}

// CredentialRevocationJob is the recurring job that revokes Vault credentials that are no
// longer being used by an active or pending session.
// The CredentialRevocationJob is not thread safe, an attempt to Run the job concurrently
// will result in an JobAlreadyRunning error.
type CredentialRevocationJob struct {
	reader db.Reader
	writer db.Writer
	kms    *kms.Kms
	logger hclog.Logger
	limit  int

	running      ua.Bool
	numCreds     int
	numProcessed int
}

// NewCredentialRevocationJob creates a new in-memory CredentialRevocationJob.
//
// WithLimit is the only supported option.
func NewCredentialRevocationJob(r db.Reader, w db.Writer, kms *kms.Kms, logger hclog.Logger, opt ...Option) (*CredentialRevocationJob, error) {
	const op = "vault.NewCredentialRevocationJob"
	switch {
	case r == nil:
		return nil, errors.New(errors.InvalidParameter, op, "missing db.Reader")
	case w == nil:
		return nil, errors.New(errors.InvalidParameter, op, "missing db.Writer")
	case kms == nil:
		return nil, errors.New(errors.InvalidParameter, op, "missing kms")
	case logger == nil:
		return nil, errors.New(errors.InvalidParameter, op, "missing logger")
	}

	opts := getOpts(opt...)
	if opts.withLimit == 0 {
		// zero signals the boundary defaults should be used.
		opts.withLimit = db.DefaultLimit
	}
	return &CredentialRevocationJob{
		reader: r,
		writer: w,
		kms:    kms,
		logger: logger,
		limit:  opts.withLimit,
	}, nil
}

// Status returns the current status of the credential revocation job.  Total is the total number
// of credentials that are set to be revoked. Completed is the number of credentials already revoked.
func (r *CredentialRevocationJob) Status() scheduler.JobStatus {
	return scheduler.JobStatus{
		Completed: r.numProcessed,
		Total:     r.numCreds,
	}
}

// Run queries the vault credential repo for credentials that need to be revoked, it then creates
// a vault client and revokes each credential.  Can not be run in parallel, if Run is invoked while
// already running an error with code JobAlreadyRunning will be returned.
func (r *CredentialRevocationJob) Run(ctx context.Context) error {
	const op = "vault.(CredentialRevocationJob).Run"
	if !r.running.CAS(r.running.Load(), true) {
		return errors.New(errors.JobAlreadyRunning, op, "job already running")
	}
	defer r.running.Store(false)

	// Verify context is not done before running
	if err := ctx.Err(); err != nil {
		return errors.Wrap(err, op)
	}

	var creds []*privateCredential
	err := r.reader.SearchWhere(ctx, &creds, "status = ?", []interface{}{RevokeCredential}, db.WithLimit(r.limit))
	if err != nil {
		return errors.Wrap(err, op)
	}

	// Set numProcessed and numTokens for status report
	r.numProcessed, r.numCreds = 0, len(creds)
	for _, c := range creds {
		// Verify context is not done before renewing next token
		if err := ctx.Err(); err != nil {
			return errors.Wrap(err, op)
		}
		if err := r.revokeCred(ctx, c); err != nil {
			r.logger.Error("error revoking credential", "credential id", c.PublicId, "error", err)
		}
		r.numProcessed++
	}

	return nil
}

func (r *CredentialRevocationJob) revokeCred(ctx context.Context, c *privateCredential) error {
	const op = "vault.(CredentialRenewalJob).revokeCred"
	databaseWrapper, err := r.kms.GetWrapper(ctx, c.ScopeId, kms.KeyPurposeDatabase)
	if err != nil {
		return errors.Wrap(err, op, errors.WithMsg("unable to get database wrapper"))
	}
	if err = c.decrypt(ctx, databaseWrapper); err != nil {
		return errors.Wrap(err, op)
	}

	vc, err := c.client()
	if err != nil {
		return errors.Wrap(err, op)
	}

	err = vc.revokeLease(c.ExternalId)
	if err != nil {
		// TODO: handle perm issue
		return errors.Wrap(err, op, errors.WithMsg("unable to revoke credential"))
	}

	cred := c.toCredential()
	query, values := cred.updateStatusQuery(RevokedCredential)
	numRows, err := r.writer.Exec(ctx, query, values)
	if err != nil {
		return errors.Wrap(err, op)
	}
	if numRows != 1 {
		return errors.New(errors.Unknown, op, "credential revoked but failed to update repo")
	}

	return nil
}

// NextRunIn determine when the next credential revocation job should run.
func (r *CredentialRevocationJob) NextRunIn() (time.Duration, error) {
	return defaultNextRunIn, nil
}

// Name is the unique name of the job.
func (r *CredentialRevocationJob) Name() string {
	return credentialRevocationJobName
}

// Description is the human readable description of the job.
func (r *CredentialRevocationJob) Description() string {
	return "Periodically revokes dynamic credentials that are no longer in use and have been set for revocation (in the revoke state)."
}
