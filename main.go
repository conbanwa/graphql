package graphql

type Filter struct {
	First        int
	Skip         int
	OrderBy      string
	OrderDescend bool //逆序
	TimeStart    int
	TimeEnd      int
}

// @version 0.1.8
