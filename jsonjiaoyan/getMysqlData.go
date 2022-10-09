package main

import (
	"database/sql"
	"encoding/json"
	_ "encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

var (
	// 定义一个全局对象db
	db *sql.DB
	//定义数据库连接的相关参数值
	//连接数据库的用户名
	userName string = "onesdev"
	//连接数据库的密码
	password string = "onesdev"
	//连接数据库的地址
	ipAddress string = "119.23.130.213"
	//连接数据库的端口号
	port int = 3306
	//连接数据库的具体数据库名称
	dbName string = "project_master"
	//连接数据库的编码格式
	charset string = "utf8"
	//数据库中为json格式的数据表以及字段 在project_master中进行查找

	/**
	"layout_card_plugin":                  {"record"},
		"third_party_backup":                  {"data"},
		"webhook_task_filter":                 {"field_filter_str", "issue_type_filter_str"},
		"account_notice":                      {"ext"},
		"activity_chart_config":               {"related_setting"},
		"activity_release":                    {"data"},
		"automation_block":                    {"data"},
		"automation_rule":                     {"filter_group"},
		"batch_task":                          {"payload", "extra", "successful", "unsuccessful", "errors_affect_count"},
		"batch_task_row":                      {"extra"},
		"card":                                {"config"},
		"common_message":                      {"ext"},
		"component":                           {"objects", "settings", "kanban_setting", "ext_settings"},
		"component_template":                  {"objects"},
		"component_view":                      {"query", "`condition`", "sort", "table_field_settings", "board_settings"},
		"component_view_user":                 {"query", "`condition`", "sort", "table_field_settings", "board_settings"},
		"field":                               {"stay_settings"},
		"file":                                {"exif"},
		"filter":                              {"query", "sort"},
		"gantt_chart":                         {"import_config"},
		"issue_type":                          {"default_configs"},
		"issue_type_layout_draft":             {"data"},
		"issue_type_notice_rule":              {"notice_time", "notice_user_domains", "filter_condition", "`condition`"},
		"layout_block":                        {"metadata"},
		"license_alter":                       {"new_info"},
		"license_change":                      {"extra"},
		"manhour_calendar":                    {"config"},
		"manhour_limit":                       {"limit_user_domains"},
		"manhour_report":                      {"config"},
		"mq_failed_message":                   {"body"},
		"object_link_type":                    {"source_condition", "target_condition"},
		"org_config":                          {"config_data"},
		"performance_config":                  {"config"},
		"product_component_view":              {"query", "`condition`", "sort", "table_field_settings", "board_settings"},
		"product_field":                       {"options"},
		"project_field":                       {"options", "statuses"},
		"project_filter":                      {"query"},
		"project_plugin":                      {"config"},
		"project_report":                      {"config"},
		"query_log":                           {"data"},
		"restore_object":                      {"data"},
		"restore_sync_object":                 {"data"},
		"team_plugin":                         {"config"},
		"team_report":                         {"config"},
		"testcase_case_recycle":               {"payload"},
		"testcase_field":                      {"options"},
		"testcase_plan_field":                 {"options", "statuses"},
		"testcase_plan_issue_tracking_config": {"column_config"},
		"testcase_report":                     {"summary"},
		"third_party_before_upgrade_data":     {"data"},
		"third_party_log":                     {"data"},
		"third_party_setting":                 {"json_config"},
		"transition":                          {"fields", "post_function"},
		"user_filter_view":                    {"query", "`condition`", "sort", "table_field_settings", "board_settings"},
		"user_reminder":                       {"ext"},
		"webhook_notice":                      {"message"},
		"workorder":                           {"snapshot"},
	*/
	//基于project_master数据库找出的json字段
	hasJsonMap = map[string][]string{
		"layout_card_plugin":                  {"record"},
		"third_party_backup":                  {"data"},
		"webhook_task_filter":                 {"field_filter_str", "issue_type_filter_str"},
		"account_notice":                      {"ext"},
		"activity_chart_config":               {"related_setting"},
		"activity_release":                    {"data"},
		"automation_block":                    {"data"},
		"automation_rule":                     {"filter_group"},
		"batch_task":                          {"payload", "extra", "successful", "unsuccessful", "errors_affect_count"},
		"batch_task_row":                      {"extra"},
		"card":                                {"config"},
		"common_message":                      {"ext"},
		"component":                           {"objects", "settings", "kanban_setting", "ext_settings"},
		"component_template":                  {"objects"},
		"component_view":                      {"query", "`condition`", "sort", "table_field_settings", "board_settings"},
		"component_view_user":                 {"query", "`condition`", "sort", "table_field_settings", "board_settings"},
		"field":                               {"stay_settings"},
		"file":                                {"exif"},
		"filter":                              {"query", "sort"},
		"gantt_chart":                         {"import_config"},
		"issue_type":                          {"default_configs"},
		"issue_type_layout_draft":             {"data"},
		"issue_type_notice_rule":              {"notice_time", "notice_user_domains", "filter_condition", "`condition`"},
		"layout_block":                        {"metadata"},
		"license_alter":                       {"new_info"},
		"license_change":                      {"extra"},
		"manhour_calendar":                    {"config"},
		"manhour_limit":                       {"limit_user_domains"},
		"manhour_report":                      {"config"},
		"mq_failed_message":                   {"body"},
		"object_link_type":                    {"source_condition", "target_condition"},
		"org_config":                          {"config_data"},
		"performance_config":                  {"config"},
		"product_component_view":              {"query", "`condition`", "sort", "table_field_settings", "board_settings"},
		"product_field":                       {"options"},
		"project_field":                       {"options", "statuses"},
		"project_filter":                      {"query"},
		"project_plugin":                      {"config"},
		"project_report":                      {"config"},
		"query_log":                           {"data"},
		"restore_object":                      {"data"},
		"restore_sync_object":                 {"data"},
		"team_plugin":                         {"config"},
		"team_report":                         {"config"},
		"testcase_case_recycle":               {"payload"},
		"testcase_field":                      {"options"},
		"testcase_plan_field":                 {"options", "statuses"},
		"testcase_plan_issue_tracking_config": {"column_config"},
		"testcase_report":                     {"summary"},
		"third_party_before_upgrade_data":     {"data"},
		"third_party_log":                     {"data"},
		"third_party_setting":                 {"json_config"},
		"transition":                          {"fields", "post_function"},
		"user_filter_view":                    {"query", "`condition`", "sort", "table_field_settings", "board_settings"},
		"user_reminder":                       {"ext"},
		"webhook_notice":                      {"message"},
		"workorder":                           {"snapshot"},
	}
)

//表字段
type Field struct {
	fieldName string
	fieldDesc string
	dataType  string
	isNull    string
	length    int
}

type FieldMy struct {
	uuidValue  string
	fieldValue string
}

/**
初始化连接数据库
*/
func initDB() (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", userName, password, ipAddress, port, dbName, charset)
	//Open打开一个driverName指定的数据库，dataSourceName指定数据源
	//不会校验用户名和密码是否正确，只会对dsn的格式进行检测
	db, err = sql.Open("mysql", dsn)
	if err != nil { //dsn格式不正确的时候会报错
		return err
	}
	//尝试与数据库连接，校验dsn是否正确
	err = db.Ping()
	if err != nil {
		fmt.Println("校验失败,err", err)
		return err
	}
	db.SetMaxIdleConns(20)

	return nil
}

//关闭数据库
func closeDB() {
	if db != nil {
		db.Close()
	}
}

//错误检查
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

//获取表属性
func TableInfo(dbName string) map[string]string {
	sqlStr := `SELECT table_name tableName,TABLE_COMMENT tableDesc
			FROM INFORMATION_SCHEMA.TABLES
			WHERE UPPER(table_type)='BASE TABLE'
			AND LOWER(table_schema) = ?
			ORDER BY table_name asc`

	var result = make(map[string]string)

	rows, err := db.Query(sqlStr, dbName)
	checkErr(err)

	for rows.Next() {
		var tableName, tableDesc string
		err = rows.Scan(&tableName, &tableDesc)
		checkErr(err)

		if len(tableDesc) == 0 {
			tableDesc = tableName
		}
		result[tableName] = tableDesc
	}
	return result
}

//获取表字段属性
func FieldInfo(dbName, tableName string) []Field {
	sqlStr := `SELECT COLUMN_NAME fName,column_comment fDesc,DATA_TYPE dataType,
						IS_NULLABLE isNull,IFNULL(CHARACTER_MAXIMUM_LENGTH,0) sLength
			FROM information_schema.columns
			WHERE table_schema = ? AND table_name = ?`

	var result []Field

	rows, err := db.Query(sqlStr, dbName, tableName)
	checkErr(err)

	for rows.Next() {
		var f Field
		err = rows.Scan(&f.fieldName, &f.fieldDesc, &f.dataType, &f.isNull, &f.length)
		checkErr(err)

		result = append(result, f)
	}
	return result
}

//获取数据库的主键字段
func GetFieldPrimay(dbName, tableName string) []string {
	sqlStr := `SELECT cu.Column_Name
		FROM  INFORMATION_SCHEMA.KEY_COLUMN_USAGE cu
WHERE CONSTRAINT_NAME = 'PRIMARY' AND cu.Table_Name = ? AND CONSTRAINT_SCHEMA=?;`

	var result []string

	rows, err := db.Query(sqlStr, tableName, dbName)
	checkErr(err)

	if rows == nil {
		return result
	}

	for rows.Next() {
		//var f Field
		var fieldName string
		err = rows.Scan(&fieldName)
		checkErr(err)

		result = append(result, fieldName)
	}
	return result
}

//SELECT cu.Column_Name
//FROM  INFORMATION_SCHEMA.`KEY_COLUMN_USAGE` cu
//WHERE CONSTRAINT_NAME = 'PRIMARY' AND cu.Table_Name = '表名' AND CONSTRAINT_SCHEMA='数据库名';

//获取表中字段属性的值
func GetValueInfo(tableName, fieldName string) []FieldMy {
	fmt.Println(tableName)

	// SELECT t1.uuid,t1.`query` from (SELECT uuid,`query` FROM   component_view_user where `query` != '' ) t1 WHERE t1.`query` != 'null';
	sqlStr := "SELECT t1.%s,t1.%s from (SELECT %s,%s FROM   %s where %s != '') t1 WHERE t1.%s != 'null';"

	//SELECT t1.%s,t1.%s from (SELECT %s,%s FROM   %s where %s != '') t1 WHERE t1.%s != 'null';
	//primay := GetFieldPrimay(tableName, dbName)
	// 获取表中字段
	infos := FieldInfo(dbName, tableName)
	biaoji := false

	for _, info := range infos {
		if info.fieldName == "uuid" {
			biaoji = true
			break
		}
	}

	//使用 regexp 正则匹配查询
	if biaoji == false {

		//sqlStr = "SELECT %s FROM  %s where %s != '';"
		sqlStr = "SELECT t1.%s from (SELECT %s FROM %s where %s != '') t1 WHERE t1.%s != 'null'"
		sqlStr = fmt.Sprintf(sqlStr, fieldName, fieldName, tableName, fieldName, fieldName)
		var result []FieldMy

		rows, err := db.Query(sqlStr)
		checkErr(err)

		for rows.Next() {
			var fieldValue sql.NullString
			err = rows.Scan(&fieldValue)
			checkErr(err)
			FieldMy := FieldMy{uuidValue: "无UUID字段列", fieldValue: fieldValue.String}

			result = append(result, FieldMy)

		}

		return result
	}
	sqlStr = fmt.Sprintf(sqlStr, "`uuid`", fieldName, "`uuid`", fieldName, tableName, fieldName, fieldName)
	var result []FieldMy

	rows, err := db.Query(sqlStr)
	checkErr(err)

	for rows.Next() {
		var fieldValue sql.NullString
		var uuidValue sql.NullString
		err = rows.Scan(&uuidValue, &fieldValue)
		checkErr(err)
		FieldMy := FieldMy{uuidValue: uuidValue.String, fieldValue: fieldValue.String}

		result = append(result, FieldMy)
	}

	return result

}

//验证字符串是否合法
func ValidJson(field string) bool {
	return json.Valid([]byte(field))
}

func getUUIDForIndex(tableName string, index int) string {

	sqlStr := "SELECT `uuid` FROM  %s LIMIT ?,?"
	sqlStr = fmt.Sprintf(sqlStr, tableName)
	var strResult string

	rows, err := db.Query(sqlStr, index-1, 1)
	checkErr(err)

	for rows.Next() {

		err = rows.Scan(&strResult)
		checkErr(err)

	}
	return strResult
}

func main() {
	err := initDB()
	if err != nil {
		fmt.Println("初始化数据库失败,err", err)
		return
	}

	for hasJsonTableName, hasJsonFieldName := range hasJsonMap {
		for _, FieldName := range hasJsonFieldName {
			fieldValues := GetValueInfo(hasJsonTableName, FieldName)
			index := 0
			for _, field := range fieldValues {
				index++
				//if field.fieldValue == "" {
				//	continue
				//}
				//if string(field.fieldValue) == "null" {
				//	continue
				//}
				//排除纯数字情况下也是json格式的问题
				panduan_one := string(field.fieldValue[0]) == "[" || string(field.fieldValue[0]) == "{"
				panduan_two := string(field.fieldValue[0]) == "{" || string(field.fieldValue[0]) == "["
				if panduan_one == false && panduan_two == false {
					fmt.Printf("解析错误数据库表名称:%s，解析错误字段名称:%s,解析错误字段值:%s \n", hasJsonTableName, FieldName, field.fieldValue)
					fmt.Printf("发生错误表的行号:%v\n", index)
					fmt.Printf("发生错误表UUID:%s\n", field.uuidValue)
					getInfoFromErrorFieldValue(hasJsonTableName, FieldName, field.fieldValue)
					return
				}

				//验证json格式是否正确
				if !ValidJson(field.fieldValue) {
					fmt.Printf("解析错误数据库表名称:%s，解析错误字段名称:%s,解析错误字段值:%s \n", hasJsonTableName, FieldName, field.fieldValue)
					fmt.Printf("发生错误的表的行号:%v\n", index)
					fmt.Printf("UUID字段的值:%s\n", field.uuidValue)
					getInfoFromErrorFieldValue(hasJsonTableName, FieldName, field.fieldValue)
					return
				}
			}
		}
	}

	//names := GetFieldPrimay("bang", "batch_task_row")
	//fmt.Printf("解析错误数据库表名称-----:%s", names)

	fmt.Printf("Json数据格式成功!!!")
	closeDB()
}

func getInfoFromErrorFieldValue(tableName string, fieldName string, fieldValue string) {

	infos := FieldInfo(dbName, tableName)
	//infos_len := len(infos)
	var str1 []string
	for _, info := range infos {
		// fmt.Println(info.fieldName)
		if strings.Contains(info.fieldName, "uuid") {
			str1 = append(str1, info.fieldName)
		}
	}
	//strings.Join(str1, ",")
	str2 := strings.Join(str1, ",")

	sqlStr := "SELECT %s FROM %s WHERE %s = '%s';"
	sqlStr = fmt.Sprintf(sqlStr, str2, tableName, fieldName, fieldValue)

	// fmt.Println(sqlStr)
	// 查询数据
	// sqlStr
	//var batch_task_uuid string
	rows := db.QueryRow(sqlStr)

	if len(str1) == 1 {
		var s1 string
		//// Scan扫描
		err := rows.Scan(&s1)
		//err := rows.Scan(&batch_task_uuid)
		if err != nil {
			fmt.Println("scan err:", err)
			return
		}
		fmt.Println("相应包含uuid字段的值")
		for _, s2 := range str1 {
			fmt.Print(s2 + " ")
		}
		fmt.Println()
		fmt.Println(s1)
	} else if len(str1) == 2 {
		var s1, s2 string
		//// Scan扫描
		err := rows.Scan(&s1, &s2)
		//err := rows.Scan(&batch_task_uuid)
		if err != nil {
			fmt.Println("scan err:", err)
			return
		}
		fmt.Println("相应包含uuid字段的值")
		for _, s2 := range str1 {
			fmt.Print(s2 + " ")
		}
		fmt.Println()
		fmt.Println(s1, s2)
	} else if len(str1) == 3 {
		var s1, s2, s3 string
		//// Scan扫描
		err := rows.Scan(&s1, &s2, &s3)
		//err := rows.Scan(&batch_task_uuid)
		if err != nil {
			fmt.Println("scan err:", err)
			return
		}
		fmt.Println("相应包含uuid字段的值")
		for _, s2 := range str1 {
			fmt.Print(s2 + " ")
		}
		fmt.Println()
		fmt.Println(s1, s2, s3)
	} else if len(str1) == 4 {
		var s1, s2, s3, s4 string
		//// Scan扫描
		err := rows.Scan(&s1, &s2, &s3, &s4)
		//err := rows.Scan(&batch_task_uuid)
		if err != nil {
			fmt.Println("scan err:", err)
			return
		}
		fmt.Println("相应包含uuid字段的值")
		for _, s2 := range str1 {
			fmt.Print(s2 + " ")
		}
		fmt.Println()
		fmt.Println(s1, s2, s3, s4)
	} else if len(str1) == 5 {
		var s1, s2, s3, s4, s5 string
		//// Scan扫描
		err := rows.Scan(&s1, &s2, &s3, &s4, &s5)
		//err := rows.Scan(&batch_task_uuid)
		if err != nil {
			fmt.Println("scan err:", err)
			return
		}
		fmt.Println("相应包含uuid字段的值")
		for _, s2 := range str1 {
			fmt.Print(s2 + " ")
		}
		fmt.Println()
		fmt.Println(s1, s2, s3, s4)
	}

}
