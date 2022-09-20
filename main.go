package graphql

type Filter struct {
	First        int
	Skip         int
	OrderBy      string
	OrderDescend bool //逆序
	// Where        map[string]string
}

// @version 0.2.1
