package constant

type SvcName string

const (
	// 认证通过后用户缓存Key
	LoginKey = "login_user"

	// Commander用于选举的状态key
	IdMachine = "id_machine"
	JobVisor  = "job_visor"
	Election  = "election"

	// svc定义
	SvcCommander SvcName = "commander"
)
