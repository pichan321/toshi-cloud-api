package handlers

import (
	"file-api/cloud"
	"file-api/structs"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetUserMetadata(c echo.Context) error {
	userClaims, ok := c.Get("user").(structs.CustomClaims)
	if !ok {
		Logger.Errorf("Error getting metadata for user: %s", c.Get("userUuid").(string))
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, userClaims)
}

func GetQuota(c echo.Context) error {
	db, err := cloud.GetPostgres()
	if err != nil {
		return ErrorHandler(c, 500, err)
	}

	userUuid := c.Get("userUuid")
	if userUuid == "" {
		return c.JSON(http.StatusBadRequest, structs.Message{Message: "A valid token is required.", Code: 400})
	}
	
	query := fmt.Sprintf(`SELECT sum(size_mb) FROM files where account_uuid = '%s'`, userUuid)
	var quota float64
	row := db.QueryRowx(query)
	row.Scan(&quota)

	defer db.Close()
	return c.JSON(http.StatusOK, fmt.Sprintf("%0.2f", quota / 1000.0))
}

// func RegisterAccount(c echo.Context) error {
// 	db, err := cloud.GetPostgres()
// 	if err != nil {
// 		return ErrorHandler(c, 500, err)
// 	}

// 	var account structs.Account
// 	err = c.Bind(&account)
// 	if err != nil {
// 		return ErrorHandler(c, 400, err)
// 	}

// 	err = utils.ValidateEmail(account.Email) 
// 	if err != nil {
// 		return ErrorHandler(c, 400, errors.New("invalid email address"))
// 	}

// 	if (account.Username == "" || account.Email == "" || account.Password == "") {
// 		return ErrorHandler(c, 400, errors.New("bad request"))
// 	}

// 	var checkAccount structs.Account
// 	accounts := db.QueryRowx(fmt.Sprintf(`select * from accounts where username = '%s' or email = '%s'`, account.Username, account.Email))
// 	accounts.StructScan(&checkAccount)

// 	if (account.Username == checkAccount.Username) || (account.Email == checkAccount.Email) {
// 		return ErrorHandler(c, 400, errors.New("Cannot register account!"))
// 	}

// 	id := uuid.New()
// 	hashedPassword := utils.HashPassword(account.Password)

// 	_, err = db.Exec(fmt.Sprintf(`insert into accounts (uuid, username, password, email, token) values ('%s', '%s','%s','%s', '%s')`,  id.String(), strings.ToLower(account.Username), string(hashedPassword), account.Email, ""))
// 	if err != nil {
// 		fmt.Printf("%v", err)
// 	}
// 	_, err = db.Exec(fmt.Sprintf(`insert into profile (uuid, link) values ('%s', '%s')`,  id.String(), ""))
// 	if err != nil {
// 		fmt.Printf("%v", err)
// 	}
// 	defer db.Close()
// 	return c.JSON(http.StatusOK, structs.Message{Code: 200})
// }

// func Login(c echo.Context) error {
// 	db, err := cloud.GetPostgres()
// 	if err != nil {
// 		return ErrorHandler(c, 500, err)
// 	}

// 	token := c.QueryParam("token")
// 	if (token != "") {
// 		row := db.QueryRowx(fmt.Sprintf(`select * from accounts where token = '%s'`, token))
// 		var account structs.Account
// 		row.StructScan(&account)
// 		if (account.Token == token) {
// 			var temp = &account
// 			temp.Password = ""
// 			return c.JSON(http.StatusOK, account)
// 		}
// 	}

// 	var account structs.Account
// 	err = c.Bind(&account)
// 	if err != nil {
// 		return ErrorHandler(c, 404, err)
// 	}

// 	row := db.QueryRowx(fmt.Sprintf(`select * from accounts where username = '%s' limit 1`, account.Username))

// 	if row == nil {
// 		return ErrorHandler(c, 404, err)
// 	}

// 	var dbAccount structs.Account
// 	err = row.StructScan(&dbAccount)
// 	if err != nil {
// 		return ErrorHandler(c, 404, err)
// 	}

// 	if (utils.HashPassword(account.Password) != dbAccount.Password || (strings.Compare(account.Username, dbAccount.Username) != 0)) {
// 		return ErrorHandler(c, 404, err)
// 	}

// 	var tempPointer = &dbAccount
// 	tempPointer.Password = ""
// 	tempPointer.Token = utils.GenerateToken()
// 	//query := fmt.Sprintf(`update accounts set token = '%s' where uuid = '%s'`, dbAccount.Token, dbAccount.Username)
// 	query := fmt.Sprintf("update accounts set token = '%s' where uuid = '%s'", dbAccount.Token, dbAccount.Uuid)
// 	_, err = db.Exec(query)
// 	if err != nil {
// 		fmt.Printf("%v", err)
// 	}

// 	defer db.Close()
// 	return c.JSON(http.StatusOK, dbAccount)
// }

// func ChangePassword(c echo.Context) error {
// 	db, err := cloud.GetPostgres()
// 	if err != nil {
// 		return ErrorHandler(c, 500, err)
// 	}
// 	var changePassword structs.ChangePassword
// 	c.Bind(&changePassword)

// 	query := fmt.Sprintf(`select count(*) from accounts where token = '%s' and password = '%s'`, changePassword.Token, utils.HashPassword(changePassword.OldPassword))

// 	row := db.QueryRowx(query)
// 	var count int
// 	row.Scan(&count)

// 	if count != 1 {
// 		return ErrorHandler(c, 400, errors.New("Cannot change password"))
// 	}

// 	query = fmt.Sprintf(`update accounts set password = '%s' where token = '%s'`, utils.HashPassword(changePassword.NewPassword), changePassword.Token)
// 	_, err = db.Exec(query)
// 	if err != nil {
// 		return ErrorHandler(c, 500, errors.New("Cannot change password"))
// 	}

// 	return c.JSON(http.StatusOK, structs.Message{Message: "Password changed!", Code: 200})
// }


// func UploadProfile(c echo.Context) error {
// 	db, err := cloud.GetPostgres()
// 	if err != nil {
// 		return ErrorHandler(c, 500, errors.New("cannot get postgres"))
// 	}
	
// 	user := c.FormValue("userUuid")
// 	file, err := c.FormFile("file")
// 	fileName := c.FormValue("fileName")
// 	bucket := utils.GetBucketUuid(5.0)
// 	fmt.Println("BUCKET PROFILE")
// 	fmt.Println(bucket)
// 	if bucket.Uuid == "" || bucket.AccessToken == "" {
// 		return ErrorHandler(c, 500, errors.New("invalid bucket"))
// 	}
	

// 	src, err := file.Open()
// 	if err != nil {
// 		return ErrorHandler(c, 500, errors.New("cannot open file"))
// 	}

// 	ctx := context.Background()
// 	project, err := cloud.GetStorj(bucket.AccessToken, ctx)
// 	if err != nil {
// 		return ErrorHandler(c, 500, errors.New("cannot open storj"))
// 	}
// 	_, err = project.EnsureBucket(context.Background(), bucket.Name)
// 	if err != nil {
// 		return ErrorHandler(c, 500, errors.New("could not ensure bucket"))
// 	}
// 	storjFileName := strings.ReplaceAll(user, "-", "_") + "_" + "profile" + "_" + fileName
// 	upload, err := project.UploadObject(ctx, bucket.Name, storjFileName, nil)
// 	if err != nil {
// 		return fmt.Errorf("could not initiate upload: %v", err)
// 	}
// 	data, err := ioutil.ReadAll(src)
// 	var oldProfile string
// 	row := db.QueryRowx(fmt.Sprintf(`select link from profile where uuid = '%s'`, user))
// 	row.Scan(&oldProfile)

// 	project.DeleteObject(ctx, bucket.Name, oldProfile)

// 	_, err = db.Exec(fmt.Sprintf(`update profile set link = '%s' where uuid = '%s'`, bucket.ShareLink + "/" + storjFileName, user))

// 	if err != nil {
// 		return fmt.Errorf("could not upload data: %v", err)
// 	}
// 	// Copy the data to the upload.
// 	buf := bytes.NewBuffer(data)
// 	_, err = io.Copy(upload, buf)
// 	if err != nil {
// 		_ = upload.Abort()
// 		return fmt.Errorf("could not upload data: %v", err)
// 	}

// 	// Commit the uploaded object.
// 	err = upload.Commit()
// 	if err != nil {
// 		return fmt.Errorf("could not commit uploaded object: %v", err)
// 	}


// 	defer db.Close()
// 	defer project.Close()
// 	defer src.Close()
// 	return c.JSON(http.StatusOK, structs.Message{Message: bucket.ShareLink + "/" + storjFileName, Code: 200})
// }

// func GetProfile(c echo.Context) error {
// 	db, _ := cloud.GetPostgres()

// 	user := c.Param("user")

// 	query := fmt.Sprintf(`select link from profile where uuid = '%s'`, user) 
// 	var link string
// 	row := db.QueryRowx(query)
// 	row.Scan(&link)


// 	return c.JSON(http.StatusOK, link + `?wrap=0`)
// }