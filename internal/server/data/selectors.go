package data

import (
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/infrahq/infra/internal/server/models"
	"github.com/infrahq/infra/uid"
)

type SelectorFunc func(db *gorm.DB) *gorm.DB

func ByID(id uid.ID) SelectorFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func ByIDs(ids []uid.ID) SelectorFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id in (?)", ids)
	}
}

func NotIDs(ids []uid.ID) SelectorFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Not(ids)
	}
}

func ByOptionalName(name string) SelectorFunc {
	return func(db *gorm.DB) *gorm.DB {
		if len(name) > 0 {
			return db.Where("name = ?", name)
		}

		return db
	}
}

func ByOptionalIDs(ids []uid.ID) SelectorFunc {
	return func(db *gorm.DB) *gorm.DB {
		if len(ids) > 0 {
			return db.Where("id in (?)", ids)
		}

		return db
	}
}

func ByName(name string) SelectorFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name = ?", name)
	}
}

func ByOptionalUniqueID(nodeID string) SelectorFunc {
	return func(db *gorm.DB) *gorm.DB {
		if len(nodeID) > 0 {
			return db.Where("unique_id = ?", nodeID)
		}

		return db
	}
}

func ByProviderID(id uid.ID) SelectorFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("provider_id = ?", id)
	}
}

func ByKeyID(key string) SelectorFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("key_id = ?", key)
	}
}

func ByOptionalSubject(polymorphicID uid.PolymorphicID) SelectorFunc {
	return func(db *gorm.DB) *gorm.DB {
		if polymorphicID == "" {
			return db
		}

		return db.Where("subject = ?", string(polymorphicID))
	}
}

func BySubject(polymorphicID uid.PolymorphicID) SelectorFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("subject = ?", string(polymorphicID))
	}
}

func ByOptionalIssuedFor(id uid.ID) SelectorFunc {
	return func(db *gorm.DB) *gorm.DB {
		if id == 0 {
			return db
		}

		return db.Where("issued_for = ?", id)
	}
}

func ByIssuedFor(id uid.ID) SelectorFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("issued_for = ?", id)
	}
}

func ByIdentityID(identityID uid.ID) SelectorFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("identity_id = ?", identityID)
	}
}

func ByUserID(userID uid.ID) SelectorFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", userID)
	}
}

func ByNotExpiredOrExtended() SelectorFunc {
	return func(db *gorm.DB) *gorm.DB {
		query := strings.Builder{}
		query.WriteString("(expires_at > ? OR expires_at = ? OR expires_at is null) AND ")
		query.WriteString("(extension_deadline > ? OR extension_deadline = ? OR extension_deadline is null)")
		return db.Where(query.String(), time.Now().UTC(), time.Time{}, time.Now().UTC(), time.Time{})
	}
}

func ByPagination(p models.Pagination) SelectorFunc {
	return func(db *gorm.DB) *gorm.DB {

		if p.Page == 0 && p.Limit == 0 {
			return db
		}
		resultsForPage := p.Limit * (p.Page - 1)
		return db.Offset(resultsForPage).Limit(p.Limit)
	}
}

func CreatedBy(id uid.ID) SelectorFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("created_by = ?", id)
	}
}

func OrderBy(order string) SelectorFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(order)
	}
}

func Limit(limit int) SelectorFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit)
	}
}

// NotCreatedBy filters out entities not created by the passed in ID
func NotCreatedBy(id uid.ID) SelectorFunc {
	return func(db *gorm.DB) *gorm.DB {
		// the created_by field is default 0 when not set by default
		return db.Where("created_by != ?", id)
	}
}

func NotName(name string) SelectorFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Not("name = ?", name)
	}
}

func NotProviderKind(kind models.ProviderKind) SelectorFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Not("kind = ?", kind)
	}
}

func ByProviderKind(kind models.ProviderKind) SelectorFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("kind = ?", kind)
	}
}

func NotPrivilege(privilege string) SelectorFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Not("privilege = ?", privilege)
	}
}

func ByOptionalIdentityGroupID(groupID uid.ID) SelectorFunc {
	return func(db *gorm.DB) *gorm.DB {
		if groupID == 0 {
			return db
		}
		return db.
			Joins("join identities_groups on identities_groups.identity_id = id").
			Where("identities_groups.group_id = ?", groupID)
	}
}
