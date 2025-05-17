package main

import (
	"context"
	"fmt"
	_ "net/http/pprof"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/bearllflee/go_shop/global"
	"github.com/bearllflee/go_shop/initialize"
	"github.com/bearllflee/go_shop/model"
	"github.com/bearllflee/go_shop/service"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Teacher struct {
	gorm.Model
	Name string
}

type Course struct {
	gorm.Model
	Name      string
	TeacherID uint
	Teacher   Teacher   `gorm:"foreignKey:TeacherID"`
	Students  []Student `gorm:"many2many:student_course;"`
}

type Class struct {
	gorm.Model
	Name     string
	Students []Student
}

type Student struct {
	gorm.Model
	Name    string
	ClassID uint
	Class   Class    `gorm:"foreignKey:ClassID"`
	Courses []Course `gorm:"many2many:student_course;"`
}

func Preload() {
	var course Course
	global.DB.First(&course, 1)
	err := global.DB.Model(&course).Association("Teacher").Find(&course.Teacher)
	if err != nil {
		panic(err)
	}
	err = global.DB.Model(&course).Association("Students").Find(&course.Students)
	if err != nil {
		panic(err)
	}
	fmt.Printf("课程名称：%s，教师名称：%s，学生数量：%d\n", course.Name, course.Teacher.Name, len(course.Students))
	for i := range course.Students {
		err = global.DB.Model(&course.Students[i]).Association("Class").Find(&course.Students[i].Class)
		if err != nil {
			panic(err)
		}
	}
	for _, student := range course.Students {
		fmt.Printf("学生名称：%s，班级名称：%s\n", student.Name, student.Class.Name)
	}
}

func QueryAssociationOneToOne() {
	// course := Course{}
	// global.DB.First(&course, 1)
	// err := global.DB.Model(&course).Association("Teacher").Find(&course.Teacher)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("课程名称：%s，教师名称：%s\n", course.Name, course.Teacher.Name)
}

func QueryAssociationOneToMany() {
	// class := Class{}
	// global.DB.First(&class, 1)
	// err := global.DB.Model(&class).Association("Students").Find(&class.Students)
	// if err != nil {
	// 	panic(err)
	// }
	// for _, student := range class.Students {
	// 	fmt.Printf("学生名称：%s，班级名称：%s\n", student.Name, class.Name)
	// }
	student := Student{}
	global.DB.First(&student, 1)
	err := global.DB.Model(&student).Association("Class").Find(&student.Class)
	if err != nil {
		panic(err)
	}
	fmt.Printf("学生名称：%s，班级名称：%s\n", student.Name, student.Class.Name)
}

func QueryAssociationManyToMany() {
	// course := Course{}
	// global.DB.First(&course, 1)
	// err := global.DB.Model(&course).Association("Students").Find(&course.Students)
	// if err != nil {
	// 	panic(err)
	// }
	// for _, student := range course.Students {
	// 	fmt.Printf("课程名称：%s，学生名称：%s\n", course.Name, student.Name)
	// }
	student := Student{}
	global.DB.First(&student, 2)
	err := global.DB.Model(&student).Association("Courses").Find(&student.Courses)
	if err != nil {
		panic(err)
	}
	for _, course := range student.Courses {
		fmt.Printf("学生名称：%s，课程名称：%s\n", student.Name, course.Name)
	}
}

func ReplaceAssociationOneToOne() {
	var course Course
	global.DB.First(&course, 4)
	err := global.DB.Model(&course).Association("Teacher").Replace(&Teacher{Model: gorm.Model{ID: 1}})
	if err != nil {
		panic(err)
	}
}

func ReplaceAssociationOneToMany() {
	var stu Student
	global.DB.First(&stu, 9)
	err := global.DB.Model(&stu).Association("Class").Replace(&Class{Model: gorm.Model{ID: 1}})
	if err != nil {
		panic(err)
	}
}

func ReplaceAssociationManyToMany() {
	student := Student{Model: gorm.Model{ID: 2}}
	err := global.DB.Model(&student).Association("Courses").Replace([]*Course{
		{Model: gorm.Model{ID: 3}, Name: "生物"},
		{Model: gorm.Model{ID: 4}, Name: "物理"},
	})
	if err != nil {
		panic(err)
	}
}

func CountAssociation() {
	// var course Course
	// global.DB.First(&course, 1)
	// count := global.DB.Model(&course).Association("Students").Count()
	// fmt.Printf("课程名称：%s，学生数量：%d\n", course.Name, count)
	var class Class
	global.DB.First(&class, 1)
	count := global.DB.Model(&class).Association("Students").Count()
	fmt.Printf("班级名称：%s，学生数量：%d\n", class.Name, count)
}

func main() {
	// Task()
	// CronTaskSimpleUse()
	initialize.MustConfig()
	go KafkaConsumer()
	KafkaProducer()
	// initialize.MustLoadZap()
	// initialize.MustInitDB()
	// initialize.AutoMigrate(global.DB)
	// ExampleJSONSerializer()
	// initialize.MustInitRedis()
	// initialize.MustRunWindowServer()
	// QueryAssociationOneToOne()
	// QueryAssociationOneToMany()
	// QueryAssociationManyToMany()
	// ReplaceAssociationOneToOne()
	// ReplaceAssociationOneToMany()
	// ReplaceAssociationManyToMany()
	// CountAssociation()
	// Preload()
	// ReflectUse()
}

// HelloConsumerGroupHandler 实现了 sarama.ConsumerGroupHandler 接口
type HelloConsumerGroupHandler struct {
}

// Setup 在消费者组创建或重新平衡时调用
func (h *HelloConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	fmt.Println("Consumer group session started")
	return nil
}

// Cleanup 在消费者组关闭或重新平衡结束时调用
func (h *HelloConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	fmt.Println("Consumer group session ended")
	return nil
}

// ConsumeClaim 处理每个分区的消息
func (h *HelloConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		fmt.Println("Received messages", string(message.Value))
		// 标记消息已处理
		session.MarkMessage(message, "")
	}
	return nil
}

func KafkaConsumer() {
	ks := service.NewKafkaService()
	ks.AddConsumerHandler("group_test", []string{model.Topic_hello}, &HelloConsumerGroupHandler{})
}

func KafkaProducer() {
	ks := service.NewKafkaService()

	for i := 1; i < 11; i++ {
		time.Sleep(1 * time.Second)
		// pmsg := &model.ProduceMsg{
		// 	ToUid:   "1",
		// 	MsgType: 1,
		// 	Content: fmt.Sprintf("test-%d", i),
		// }
		// pmsgBytes, err := json.Marshal(pmsg)
		// if err != nil {
		// 	fmt.Println(err)
		// 	continue
		// }
		// pmsgStr := base64.StdEncoding.EncodeToString(pmsgBytes)
		// bmsg := model.BaseMsg{
		// 	MsgID:   model.Msg_hello_id,
		// 	MsgData: pmsgStr,
		// }
		// jsonBytes, _ := json.Marshal(bmsg)
		ks.ProduceMsg(model.Topic_hello, fmt.Sprintf("hello-%d", i))
	}
}

func init() {
	schema.RegisterSerializer("csv", CSVSerializer{})
}

// 定义一个包含指针字段的结构体
type MyStruct struct {
	Data *[]string // 注意：这是一个指向 string slice 的指针
}

func ReflectUse() {
	// var myStruct MyStruct

	// 创建一个指向 string slice 的指针
	var slice *[]string
	// myStruct.Data = &slice
	slice = new([]string)
	// 使用反射设置字段值
	field := reflect.ValueOf(slice).Elem()
	field.Set(reflect.ValueOf([]string{"apple", "apple", "cherry"}))
	fmt.Println(slice)
}

// 自定义 CSV 序列化器
type CSVSerializer struct{}

func (CSVSerializer) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue any) error {
	fmt.Println("dst.Type():", dst.Type())
	if dbValue == nil {
		return nil
	}
	var str string
	switch v := dbValue.(type) {
	case string:
		str = v
	case []byte:
		str = string(v)
	default:
		return fmt.Errorf("不支持的类型: %T", dbValue)
	}
	values := strings.Split(str, ",")
	val := field.ReflectValueOf(ctx, dst)
	// dst.Elem().FieldByName(field.Name).Set(reflect.ValueOf(values))
	val.Set(reflect.ValueOf(values))
	return nil
}

// 实现序列化方法
func (CSVSerializer) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
	if v, ok := fieldValue.([]string); ok {
		return strings.Join(v, ","), nil
	}
	return nil, fmt.Errorf("不支持的类型: %T", fieldValue)
}

func ExampleJSONSerializer() {
	course := &model.Subject{
		Name:     "Go高级编程111",
		Tags:     []string{"Golang111", "后端", "高级"},
		Syllabus: []string{"并发编程", "网络编程", "底层原理"},
		Properties: map[string]any{
			"难度级别": "高级",
			"适合人群": "有Go基础的开发者",
			"预计学时": 48,
		},
	}

	// 创建记录
	global.DB.Create(course)

	// 查询记录
	var result model.Subject
	global.DB.First(&result, course.ID)

	fmt.Printf("课程名称: %s\n", result.Name)
	fmt.Printf("课程标签: %v\n", result.Tags)
	fmt.Printf("课程大纲: %v\n", result.Syllabus)
	fmt.Printf("课程属性: %v\n", result.Properties)
}

func Task() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	wg := sync.WaitGroup{}
	wg.Add(1)
	quit := make(chan struct{})
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ticker.C:
				fmt.Println("定时任务执行：", time.Now().Format("2006-01-02 15:04:05"))
			case <-quit:
				fmt.Println("Quit")
				return
			}
		}
	}()
	time.Sleep(9 * time.Second)
	close(quit)
	// quit <- struct{}{}
	wg.Wait()
}

func CronTaskSimpleUse() {
	c := cron.New(cron.WithSeconds())
	c.Start()
	id, err := c.AddFunc("*/2 * * * * *", func() {
		fmt.Println("定时任务执行：", time.Now().Format("2006-01-02 15:04:05"))
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(id)
	time.Sleep(9 * time.Second)
	c.Stop()
}
