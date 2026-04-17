package errx

import "errors"

var (
	ErrDecryptFailed    = errors.New("decrypt failed")
	ErrEncryptFailed    = errors.New("encrypt failed")
	ErrHashingFailed    = errors.New("hashing failed")
	ErrPasswordMismatch = errors.New("password mismatch")
)

var (
	TooManyOTPError = errors.New("Please wait 1 minute before requesting another OTP")
)

// ----------------------err code --------------------
const (
	UnknownError ApiError = "UNKNOWN"
)
const (
	InvalidCredentialsCode ApiError = "INVALID_CREDENTIALS"
	UserAlreadyExistsCode  ApiError = "USER_ALREADY_EXISTS"
	PermissionDeniedCode   ApiError = "PERMISSION_DENIED"
	UserNotFoundCode       ApiError = "USER_NOT_FOUND"
)
const (
	UnknownErrorCode        ApiError = "UNKNOWN_ERROR"
	InternalServerErrorCode ApiError = "INTERNAL_SERVER_ERROR"
	BadRequestCode          ApiError = "BAD_REQUEST"
	UnauthorizedCode        ApiError = "UNAUTHORIZED"
	ForbiddenCode           ApiError = "FORBIDDEN"
	NotFoundCode            ApiError = "NOT_FOUND"
	ConflictCode            ApiError = "CONFLICT"
	TooManyRequestsCode     ApiError = "TOO_MANY_REQUESTS"
	ValidationErrorCode     ApiError = "VALIDATION_ERROR"
	UnprocessableEntityCode ApiError = "UNPROCESSABLE_ENTITY"
	MethodNotAllowedCode    ApiError = "METHOD_NOT_ALLOWED"
)
const (
	DBErrorCode               ApiError = "DATABASE_ERROR"
	DBUniqueViolationCode     ApiError = "DB_UNIQUE_VIOLATION"
	DBForeignKeyViolationCode ApiError = "DB_FOREIGN_KEY_VIOLATION"
	DBNotFoundCode            ApiError = "DB_NOT_FOUND"
	DBTimeoutCode             ApiError = "DB_TIMEOUT"
	DBConnectionFailedCode    ApiError = "DB_CONNECTION_FAILED"
)
const (
	OTPNotFoundCode         ApiError = "OTP_NOT_FOUND"
	InvalidOTPCode          ApiError = "INVALID_OTP"
	ExpiredOTPCode          ApiError = "EXPIRED_OTP"
	TooManyOTPRequestsCode  ApiError = "TOO_MANY_OTP_REQUESTS"
	OTPGenerationFailedCode ApiError = "OTP_GENERATION_FAILED"
	EmailSendFailedCode     ApiError = "EMAIL_SEND_FAILED"
)
