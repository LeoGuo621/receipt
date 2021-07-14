package routes

import (
	"errors"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"receipt/dbs"
	"receipt/eths"
	"receipt/utils"
)

var session *sessions.CookieStore

func init()  {
	session = sessions.NewCookieStore([]byte("secret"))
}
// handler

// response
func ResponseData(c echo.Context, resp *utils.Resp) {
	resp.Msg = utils.RecodeText(resp.Recode)
	_ = c.JSON(http.StatusOK, resp)
}

// GET: /ping
func PingHandler(c echo.Context) error {
	var resp utils.Resp
	resp.Recode = utils.RecodeOk
	defer ResponseData(c, &resp)
	return nil
}

// POST: /register
func Register(c echo.Context) error {
	//1. new Resp instance
	resp := utils.Resp{
		Recode: utils.RecodeOk,
	}
	defer ResponseData(c, &resp)
	//2. parse request
	user := &dbs.User{}
	if err := c.Bind(user); err != nil {
		fmt.Println("Failed to bind user: ", err)
		resp.Recode = utils.RecodeParamErr
		return err
	}
	//3. create user address
	addr, err := eths.NewAccount(user.Password)
	if err != nil {
		fmt.Println("Failed to NewAccount: ", err)
		resp.Recode = utils.RecodeEthErr
		return err
	}
	user.Address = addr
	//4. save data to mysql
	if err := user.AddUser(); err != nil {
		fmt.Println("Failed to user.AddUser(): ", err)
		resp.Recode = utils.RecodeDBErr
		return err
	}
	resp.Data = addr
	return nil
}

// POST: /login
func Login(c echo.Context) error {
	//1. new Resp instance
	resp := utils.Resp{
		Recode: utils.RecodeOk,
	}
	defer ResponseData(c, &resp)
	//2. parse request
	user := &dbs.User{}
	if err := c.Bind(user); err != nil {
		fmt.Println("Failed to bind user: ", err)
		resp.Recode = utils.RecodeParamErr
		return err
	}
	//3. query whether user exists
	ok, err := user.QueryUser()
	if err != nil {
		fmt.Println("Failed to user.QueryUser(): ", err)
		resp.Recode = utils.RecodeDBErr
		return err
	}
	if !ok {
		fmt.Println("Username or password error!")
		resp.Recode = utils.RecodeLoginErr
	}
	//4. session
	sess, _ := session.Get(c.Request(), "session")
	sess.Options = &sessions.Options{
		Path: "/",
		HttpOnly: true,
	}
	sess.Values["address"] = user.Address
	// actually, it is not recommended to save user password
	// here is just for the convenience of development
	sess.Values["password"] = user.Password
	_ = sess.Save(c.Request(), c.Response())

	resp.Data = user.Address
	return nil
}

// GET: /session
func Session(c echo.Context) error {
	var resp utils.Resp
	resp.Recode = utils.RecodeOk
	defer ResponseData(c, &resp)

	sess, err := session.Get(c.Request(), "session")
	if err != nil {
		resp.Recode = utils.RecodeLoginErr
		return err
	}
	address := sess.Values["address"]
	if address == "" {
		fmt.Println("Failed to get session, user is nil!")
		resp.Recode = utils.RecodeLoginErr
		return err
	}
	return nil
}

// Post: /content
func Upload(c echo.Context) error {
	//1. new Resp instance
	var resp utils.Resp
	resp.Recode = utils.RecodeOk
	defer ResponseData(c, &resp)

	//2. parse content
	content := &dbs.Content{}
	// FormFile returns the multipart form file for the provided name.
	// frontend and backend need to negotiate the name, here "fileName" is used temporarily.
	fh, err := c.FormFile("fileName")
	if err != nil {
		fmt.Println("Failed to FormFile: ", err)
		resp.Recode = utils.RecodeParamErr
		return err
	}
	srcFile, err := fh.Open()
	if err != nil {
		fmt.Println("Failed to fh.Open(): ", err)
	}
	defer srcFile.Close()
	//2.2 get tokenid
	tokenid := utils.NewTokenID()
	content.TokenID = fmt.Sprintf("%d", tokenid)
	filename := fmt.Sprintf("static/contents/%s.jpg", content.TokenID)
	content.ContentPath = fmt.Sprintf("/contents/%s.jpg", content.TokenID)
	dst, err := os.Create(filename)
	if err != nil {
		fmt.Println("Failed to create file: ", err, content.ContentPath)
	}
	defer dst.Close()

	//2.3 generate hash
	cData := make([]byte, fh.Size)
	n, err := srcFile.Read(cData)
	if err != nil && fh.Size != int64(n){
		resp.Recode = utils.RecodeSysErr
		return err
	}
	hash := eths.KeccakHash(cData)
	content.ContentHash = fmt.Sprintf("%x", hash)
	_, _ = dst.Write(cData)
	content.Title = fh.Filename

	//3. get account address from session
	sess, _ := session.Get(c.Request(), "session")
	content.Address, _ = sess.Values["address"].(string)
	// get user password from session. However, it is not safe.
	pw, ok := sess.Values["password"].(string)
	if !ok || content.Address == "" || pw == "" {
		resp.Recode = utils.RecodeLoginErr
		return errors.New("no session")
	}

	//4. save data to database
	err = content.AddContent()
	if err != nil {
		resp.Recode = utils.RecodeDBErr
		return err
	}

	//TODO: 5. interact ethereum

	return nil
}

const PageMaxContent = 5
// GET: /content
func GetContents(c echo.Context) error {
	//1. new Resp instance
	var resp utils.Resp
	resp.Recode = utils.RecodeOk
	defer ResponseData(c, &resp)

	//2. get user address from session
	sess, err := session.Get(c.Request(), "session")
	if err != nil {
		fmt.Println("failed to get session")
		resp.Recode = utils.RecodeLoginErr
		return err
	}
	address, ok := sess.Values["address"].(string)
	if address == "" || !ok {
		resp.Recode = utils.RecodeLoginErr
		return errors.New("no session")
	}
	//3. query database
	contents, err := dbs.QueryContent(address)
	if err != nil {
		resp.Recode = utils.RecodeDBErr
		return err
	}
	//4. arrange response data
	mapResp := make(map[string]interface{})
	mapResp["total_page"] = len(contents) /PageMaxContent + 1
	mapResp["current_page"] = 1
	mapResp["contents"] = contents
	resp.Data = mapResp
	return nil
}
