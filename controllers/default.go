package controllers

type MainController struct {
	Common
}

func (c *MainController) Get() {
	c.Redirect("/admin/login/login",302)
	c.StopRun()

}
