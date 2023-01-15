package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/graphql-go/graphql"
)

type Student struct {
	StuName string `json:"stu_name"`
	StuAge  int    `json:"stu_age"`
}

var FieldStudentType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Student",
	Fields: graphql.Fields{
		"stu_name": &graphql.Field{Type: graphql.String},
		"stu_age": &graphql.Field{Type: graphql.String},
	},
})

//班级
type ClassType struct {
	ClassName string `json:"class_name"`
	ClassNum  int    `json:"class_num"` //级别人数
	Students   []Student  `json:"students"` //学生

}

var FieldClassType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ClassType",
	Fields: graphql.Fields{
		"class_name": &graphql.Field{Type: graphql.String},
		"class_num": &graphql.Field{Type: graphql.String},
		"students" : &graphql.Field{
			Type: graphql.NewList(FieldStudentType),
		},
	},
})


type School struct {
	SchoolName string `json:"school_name"`
	SchoolAge  int     `json:"sch_age"`
	ClassList  []ClassType  `json:"class_list"`
}

var FieldSchoolType = graphql.NewObject(graphql.ObjectConfig{
	Name: "School",
	Fields: graphql.Fields{
		"school_name": &graphql.Field{Type: graphql.String},
		"sch_age": &graphql.Field{Type: graphql.String},
		"class_list": &graphql.Field{
			Type: graphql.NewList(FieldClassType),
		},
	},
})

type Bar struct {
	Name string
}

var FieldBarType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Bar",
	Fields: graphql.Fields{
		"name": &graphql.Field{Type: graphql.String},
	},
})


// QueryType fields: `concurrentFieldFoo` and `concurrentFieldBar` are resolved
// concurrently because they belong to the same field-level and their `Resolve`
// function returns a function (thunk).
var QueryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"FieldSchool": &graphql.Field{
			Type: graphql.NewList(FieldSchoolType), //表示返回数组类型
			//Type: FieldSchoolType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {

				schools := make([]*School,0)
				var foo1 = School{
					SchoolName: "xx school",
					SchoolAge: 14,
					ClassList: []ClassType{
						{
							ClassName: "中一班",
							ClassNum:  18,
							Students: []Student{
								{
									StuName: "guanzilin",
									StuAge:  5,
								},

							},
						},
					},
				}
				var foo2 = School{
					SchoolName: "YY school",
					SchoolAge: 15,
					ClassList: []ClassType{
						{
							ClassName: "中一班",
							ClassNum:  19,
							Students: []Student{
								{
									StuName: "guanzilin22",
									StuAge:  6,
								},
							},
						},
					},
				}
                schools = append(schools,&foo1)
                schools = append(schools,&foo2)//如果返回多个的话 也就是数组
				return func() (interface{}, error) {
					return schools, nil
				}, nil
			},
		},
		"name": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var NameStr = "AAAA"
				return NameStr,nil
			},
		},
		//"concurrentFieldBar": &graphql.Field{
		//	Type: FieldBarType,
		//	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		//		type result struct {
		//			data interface{}
		//			err  error
		//		}
		//		ch := make(chan *result, 1)
		//		go func() {
		//			defer close(ch)
		//			bar := &Bar{Name: "Bar's name"}
		//			ch <- &result{data: bar, err: nil}
		//		}()
		//		return func() (interface{}, error) {
		//			r := <-ch
		//			return r.data, r.err
		//		}, nil
		//	},
		//},
	},
})

func main() {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: QueryType,
	})
	if err != nil {
		log.Fatal(err)
	}
	//query := `
	//	query {
	//		concurrentFieldFoo {
	//			name
	//		}
	//		concurrentFieldBar {
	//			name
	//		}
	//	}
	//`

	query := `
		query {
			FieldSchool {
				school_name
                sch_age
                class_list {
	                class_name
	                class_num
                    students {
                       stu_name
                       stu_age
                    }
                }
			}
            name
		}
	`
	result := graphql.Do(graphql.Params{
		RequestString: query,
		Schema:        schema,
	})
	b, err := json.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", b)
	/*
		{
		  "data": {
		    "concurrentFieldBar": {
		      "name": "Bar's name"
		    },
		    "concurrentFieldFoo": {
		      "name": "Foo's name"
		    }
		  }
		}
	*/
}
