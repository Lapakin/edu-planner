package query

const MaxParams = 65535

const (
	ID         = "id"
	IsDeleted  = "is_deleted"
	IsActive   = "is_active"
	CreatedAt  = "created_at"
	ModifiedAt = "modified_at"
)

type Operator string

const (
	And Operator = "AND"
	Or  Operator = "OR"
)

type ConditionOperator string

const (
	Equals             ConditionOperator = "="
	NotEquals          ConditionOperator = "!="
	GreaterThan        ConditionOperator = ">"
	GreaterThanOrEqual ConditionOperator = ">="
	LessThan           ConditionOperator = "<"
	LessThanOrEqual    ConditionOperator = "<="
	Like               ConditionOperator = "LIKE"
	ILike              ConditionOperator = "ILIKE"
	In                 ConditionOperator = "IN"
	NotIn              ConditionOperator = "NOT IN"
	IsNull             ConditionOperator = "IS NULL"
	IsNotNull          ConditionOperator = "IS NOT NULL"
	Contains           ConditionOperator = "CONTAINS"
)

type Condition struct {
	Name     string
	Column   string
	Operator ConditionOperator
}

type Conditions []*Condition

type Where struct {
	Operator   Operator
	Conditions Conditions
}

type Wheres []*Where
