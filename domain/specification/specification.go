package specification

import "ddd-bottomup/domain/entity"

// CircleSpecification サークル専用のSpecificationインターフェース
type CircleSpecification interface {
	IsSatisfiedBy(circle *entity.Circle) bool
}