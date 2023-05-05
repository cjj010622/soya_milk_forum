package errs

//成功返回
const OK uint32 = 200

/**(前3位代表业务,后三位代表具体功能)**/

//全局错误码
const SERVER_COMMON_ERROR uint32 = 100001
const REUQEST_PARAM_ERROR uint32 = 100002
const TOKEN_EXPIRE_ERROR uint32 = 100003
const TOKEN_GENERATE_ERROR uint32 = 100004
const DB_ERROR uint32 = 100005
const DB_UPDATE_AFFECTED_ZERO_ERROR uint32 = 100006

//用户模块
const TELEPHONE_ALREADY_REGIST uint32 = 100007
const CREATE_TOKEN_ERR uint32 = 100008
const USER_NOT_EXIST uint32 = 100009
const PASSWORD_ERR uint32 = 100010
