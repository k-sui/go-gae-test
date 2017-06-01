package hello

import (

	"net/http"
	"html/template"
//	"strconv"
//	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/mail"
//	"appengine_internal/datastore"
	"google.golang.org/appengine/datastore"
	"regexp"

	"golang.org/x/net/context"

)

var tpl *template.Template

const PASSWORD = "password"

type User struct {
	Id int64
	Name string
}

type ApplyContent struct{
	Name string
	Department string
	Phone string
	Email string
	First string
	Second string
	Third string
	Fourth string
	Fifth string
	Password string

	NameCheck string
	DepartmentCheck string
	PhoneCheck string
	EmailCheck string
	PasswordCheck string

}

func init() {
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	http.HandleFunc("/apply.html", apply_handler)
	http.HandleFunc("/confirm.html", confirm_handler)
	http.HandleFunc("/result.html", result_handler)
	http.HandleFunc("/manage.html", manage_handler)
	http.HandleFunc("/drawlot.html", drawlot_handler)
	http.HandleFunc("/addlot.html", addlot_handler)
	http.HandleFunc("/", handler)

}

func handler(w http.ResponseWriter, r *http.Request) {

	var user User
	user.Name = "view controller test message"
	user.Id = 100
	tpl = template.Must(template.ParseFiles("index.html"))
	tpl.Execute(w,user)
}

func apply_handler(w http.ResponseWriter, r *http.Request) {


	var applyContent ApplyContent
	check := r.FormValue("formcheck")

//	fmt.Println("出力チェック")
/*	if check != "true"{
		tpl = template.Must(template.ParseFiles("apply.html"))
		tpl.Execute(w,applyContent)
		return
	}
*/
	c := appengine.NewContext(r)
	log.Debugf(c,check)


	applyContent.Name = r.FormValue("name")
	applyContent.Department = r.FormValue("department")
	applyContent.Phone = r.FormValue("phone")
	applyContent.Email = r.FormValue("email")
	applyContent.First = r.FormValue("first")
	applyContent.Second = r.FormValue("second")
	applyContent.Third = r.FormValue("third")

	// 全て未入力(＝初アクセス)ではなく、入力されていない情報があれば警告を出す
	if !checkApplyNoInput(&applyContent){
		checkApplyInput(&applyContent)
	}

	log.Debugf(c,applyContent.First)
	log.Debugf(c,applyContent.Second)
	log.Debugf(c,applyContent.Third)

	tpl = template.Must(template.ParseFiles("apply.html"))
	tpl.Execute(w,applyContent)

}

func confirm_handler(w http.ResponseWriter, r *http.Request){

	check := r.FormValue("formcheck")

	//	fmt.Println("出力チェック")
	if check != "true"{
//		tpl = template.Must(template.ParseFiles("apply.html"))
//		tpl.Execute(w,"")
		apply_handler(w,r)
		return
	}

	c := appengine.NewContext(r)
	log.Debugf(c,check)
	var applyContent ApplyContent

	applyContent.Name = r.FormValue("name")
	applyContent.Department = r.FormValue("department")
	applyContent.Phone = r.FormValue("phone")
	applyContent.Email = r.FormValue("email")
	applyContent.First = r.FormValue("first")
	applyContent.Second = r.FormValue("second")
	applyContent.Third = r.FormValue("third")

	// 入力してない項目がある、またはメールアドレスが違えばapply.htmlに戻す
	if !checkApplyInput(&applyContent){
		apply_handler(w,r)
		return
	}

	//ここにアクセスした時にpasswordが入力されていれば、パスワードを間違えているので、チェックにerrを入れる
	if r.FormValue("password")!=""{
		applyContent.PasswordCheck = "err"
	}

	tpl = template.Must(template.ParseFiles("confirm.html"))
	tpl.Execute(w,applyContent)

}

func result_handler(w http.ResponseWriter, r *http.Request){

	check := r.FormValue("formcheck")

	//	fmt.Println("出力チェック")
	if check != "true"{
//		tpl = template.Must(template.ParseFiles("apply.html"))
//		tpl.Execute(w,"")
		apply_handler(w,r)
		return
	}

	c := appengine.NewContext(r)

	var applyContent ApplyContent


	applyContent.Name = r.FormValue("name")
	applyContent.Department = r.FormValue("department")
	applyContent.Phone = r.FormValue("phone")
	applyContent.Email = r.FormValue("email")
	applyContent.First = r.FormValue("first")
	applyContent.Second = r.FormValue("second")
	applyContent.Third = r.FormValue("third")
	pwd := r.FormValue("password")

	if pwd!=PASSWORD{
//		applyContent.Password="true"
		confirm_handler(w,r)
		return
	}

	var applies []ApplyContent
	query := datastore.NewQuery("apply").Filter("Email=",applyContent.Email)
	keys,err := query.GetAll(c,&applies)
	//TODO keyが2個以上見つかった場合の処理
	// DataStoreへのput
	var key *datastore.Key
	if len(keys) >0{
		key = keys[0]
	}else {
		key = datastore.NewIncompleteKey(c, "apply", nil)
	}
	key, err = datastore.Put(c, key, &applyContent)
	if err != nil{
		log.Debugf(c,"エラー")
		tpl = template.Must(template.ParseFiles("apply.html"))
		tpl.Execute(w,"")
		return
	}

	//DataStoreからのGet
	var applyResult ApplyContent
	err = datastore.Get(c, key, &applyResult)

	if err != nil{
		log.Debugf(c,"エラー")
		tpl = template.Must(template.ParseFiles("apply.html"))
		tpl.Execute(w,"")
		return
	}

	sendResultMail(c, applyResult)

	tpl = template.Must(template.ParseFiles("result.html"))
	tpl.Execute(w,applyResult)

}


func manage_handler(w http.ResponseWriter, r *http.Request) {

	var user User
	user.Name = "view controller test message"
	user.Id = 100
	tpl = template.Must(template.ParseFiles("manage.html"))
	tpl.Execute(w,user)
}

func drawlot_handler(w http.ResponseWriter, r *http.Request) {

	var user User
	user.Name = "view controller test message"
	user.Id = 100
	tpl = template.Must(template.ParseFiles("lot.html"))
	tpl.Execute(w,user)
}

func addlot_handler(w http.ResponseWriter, r *http.Request) {

	var user User
	user.Name = "view controller test message"
	user.Id = 100
	tpl = template.Must(template.ParseFiles("addlot.html"))
	tpl.Execute(w,user)
}




// メールアドレスが正しいか確認
func checkMailAddress (mail string) bool {

//	reg := "/^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\\.[a-zA-Z0-9-]+)*$/"
//	re1 := regexp.MustCompile(reg).Match([]byte(mail))
//	if !re1{return false}

	re2,_ :=regexp.MatchString(".+@gmail\\.com",mail)


	return re2
}

// ApplyContentに空文字がないことを確認する。空文字があればチェックにerrを入れる
func checkApplyInput(applyContent *ApplyContent) bool{

	var result bool
	result = true

	if applyContent.Name == ""{
		result = false
		applyContent.NameCheck = "err"
	}
	if applyContent.Department == ""{
		result = false
		applyContent.DepartmentCheck = "err"
	}
	if applyContent.Phone == ""{
		result = false
		applyContent.PhoneCheck = "err"
	}
	//メールアドレスのみチェック
	if applyContent.Email == "" || !checkMailAddress(applyContent.Email){
		result = false
		applyContent.EmailCheck = "err"
	}
	if applyContent.First == ""{
		result = false
	}
	if applyContent.Second == ""{
		result = false
	}
	if applyContent.Third == ""{
		result = false
	}

	return result
}

// ApplyContentが全て空文字であることを確認する
func checkApplyNoInput(applyContent *ApplyContent) bool {


	if applyContent.Name != ""{return false}
	if applyContent.Department != ""{return false}
	if applyContent.Phone != ""{return false}
	if applyContent.Email != ""{return false}
	if applyContent.First != ""{return false}
	if applyContent.Second != ""{return false}
	if applyContent.Third != ""{return false}

	return true

}

//申込内容をメール
func sendResultMail(c context.Context, applyContent ApplyContent) bool{


	sender := "sample@sample.com"
	to := []string{applyContent.Email}
	subject := "申込受付確認メール"

	body := "以下の内容で申し込みを受け付けました\n\n"
	body += "名前　　　　：" + applyContent.Name + "\n"
	body += "所属部署　　：" + applyContent.Department + "\n"
	body += "メール　　　：" + applyContent.Email + "\n"
	body += "電話番号　　：" + applyContent.Phone + "\n"
	body += "第1希望　   ：" + applyContent.First + "\n"
	body += "第2希望　   ：" + applyContent.Second + "\n"
	body += "第3希望　   ：" + applyContent.Third + "\n"

	msg := &mail.Message{
		Sender: sender,
		To: to,
		Subject: subject,
		Body: body,
		HTMLBody: "",
	}

	err := mail.Send(c, msg)

	if err != nil{
		log.Debugf(c,"エラー")
		return false
	}

	return true

}