package passwordRecovery

type PasswordRecoveryRequest struct {
	AttemptId   string `json:"attemptId"`
	NewPassword string `json:"newPassword"`
}
