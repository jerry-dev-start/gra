package validate

// 预定义规则集 — 按需取用，可跨 struct 复用

var (
	// 通用
	IDRules = Rules{"ID": {Required("ID不能为空")}}

	// 分页
	PageRules = Rules{
		"Page": {Required("页码不能为空"), Ge("1")},
		"Size": {Required("每页条数不能为空"), Ge("1"), Le("100")},
	}

	// 用户模块
	LoginRules = Rules{
		"Username": {Required("用户名不能为空")},
		"Password": {Required("密码不能为空")},
	}

	CreateUserRules = Rules{
		"Username": {Required("用户名不能为空"), Ge("3"), Le("64")},
		"Password": {Required("密码不能为空"), Ge("6")},
	}
)
