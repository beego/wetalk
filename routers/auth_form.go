package routers

type RegisterForm struct {
	UserName   string `form:"username" valid:"Required;AlphaDash;MinSize(5);MaxSize(30)"`
	Email      string `form:"email" valid:"Required;Email;MaxSize(80)"`
	Password   string `form:"password" valid:"Required;MinSize(4);MaxSize(30)"`
	PasswordRe string `form:"passwordre" valid:"Required;MinSize(4);MaxSize(30)"`
}

type LoginForm struct {
	UserName string `form:"username" valid:"Required"`
	Password string `form:"password" valid:"Required"`
}
