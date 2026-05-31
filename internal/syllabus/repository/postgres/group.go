package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/Lapakin/edu-planner/internal/domain"

	q "github.com/Lapakin/edu-planner/internal/adapter/db/postgres/query"
	f "github.com/Lapakin/edu-planner/internal/app/filter"
)

type GroupRepository struct {
	db sqlx.ExtContext
}

func NewGroupRepository(db sqlx.ExtContext) *GroupRepository {
	return &GroupRepository{
		db: db,
	}
}

func (r *GroupRepository) CreateGroups(ctx context.Context, groups domain.Groups) error {
	builder := q.Insert(ctx).
		Into("education_group").
		Columns(`
			specialty_id,
			name,
			short_name,
			student_count,
			is_contract,
			is_splitting,
			education_start,
			education_end,
			created_at
		`)

	for _, g := range groups {
		builder = builder.Values(
			g.SpecialtyID,
			g.Name,
			g.ShortName,
			g.StudentCount,
			g.IsContract,
			g.IsSplitting,
			g.EducationStart,
			g.EducationEnd,
			g.CreatedAt,
		)
	}

	return builder.
		ReturningID().
		QueryxContext(r.db, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&groups[*i].ID); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r *GroupRepository) GetGroupByID(ctx context.Context, id uint64) (*domain.Group, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			g.id,
			gp.academic_year_id,
			g.specialty_id,
			g.name,
			g.short_name,
			g.student_count,
			g.is_contract,
			g.is_splitting,
			g.education_start,
			g.education_end,
			g.created_at,
			g.modified_at
		`).
		From("education_group g").
		InnerJoin("academic_year_to_education_group gp ON g.id = gp.group_id").
		WhereID(id).
		IsDeleted(false).
		ToSQL()
	if err != nil {
		return nil, err
	}

	group := &domain.Group{}
	if err = sqlx.GetContext(ctx, r.db, group, query, args...); err != nil {
		return nil, err
	}

	return group, nil
}

func (r *GroupRepository) FetchGroups(ctx context.Context, filters f.Filters) (domain.Groups, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			g.id,
			gp.academic_year_id,
			g.specialty_id,
			g.name,
			g.short_name,
			g.student_count,
			g.is_contract,
			g.is_splitting,
			g.education_start,
			g.education_end,
			g.created_at,
			g.modified_at
		`).
		From("education_group g").
		InnerJoin("academic_year_to_education_group gp ON g.id = gp.group_id").
		ApplyQueryFilters(filters, q.Wheres{
			{
				Operator: q.And,
				Conditions: q.Conditions{
					{Name: domain.IDsParam, Column: "g.id", Operator: q.Equals},
					{Name: domain.AcademicYearIDParam, Column: "gp.academic_year_id", Operator: q.Equals},
					{Name: domain.SpecialtyIDParam, Column: "g.specialty_id", Operator: q.Equals},
				},
			},
		}).
		IsDeleted(false).
		OrderBy("g.id").
		ToSQL()
	if err != nil {
		return nil, err
	}

	groups := make(domain.Groups, 0)
	if err = sqlx.SelectContext(ctx, r.db, &groups, query, args...); err != nil {
		return nil, err
	}

	return groups, nil
}

func (r *GroupRepository) UpdateGroups(ctx context.Context, groups domain.Groups) error {
	builder := q.BatchUpdate(ctx).
		Table("education_group").
		PrimaryKey("id").
		NumberOfRows(len(groups)).
		AddBigintColumn("id").
		AddBigintColumn("specialty_id").
		AddVarCharColumn("name").
		AddVarCharColumn("short_name").
		AddIntColumn("student_count").
		AddBoolColumn("is_contract").
		AddBoolColumn("is_splitting").
		AddTimeStampColumn("education_start").
		AddTimeStampColumn("education_end").
		AddTimeStampColumn("modified_at")

	args := make([]any, 0, 10*len(groups))
	for _, g := range groups {
		args = append(args,
			g.ID,
			g.SpecialtyID,
			g.Name,
			g.ShortName,
			g.StudentCount,
			g.IsContract,
			g.IsSplitting,
			g.EducationStart,
			g.EducationEnd,
			g.ModifiedAt,
		)
	}

	return builder.
		ReturningCreatedAt().
		QueryxContext(r.db, args, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&groups[*i].CreatedAt); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r *GroupRepository) DeleteGroups(ctx context.Context, ids []uint64) error {
	query, args, err := q.SoftDelete(ctx).
		Table("education_group").
		WhereID(ids).
		ToSQL()
	if err != nil {
		return err
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *GroupRepository) AttachGroupsToAcademicYear(ctx context.Context, toAttach map[uint64]domain.Groups) error {
	builder := q.Insert(ctx).
		Into("academic_year_to_education_group").
		Columns(`
			academic_year_id,
			group_id
		`)

	for academicYearID, groups := range toAttach {
		for _, g := range groups {
			builder = builder.Values(
				academicYearID,
				g.ID,
			)
		}
	}

	query, args, err := builder.ToSQL()
	if err != nil {
		return err
	}

	if _, err = r.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (r *GroupRepository) CreateGroupSemesters(ctx context.Context, groupSemesters domain.GroupSemesters) error {
	builder := q.Insert(ctx).
		Into("education_group_semester").
		Columns(`
			group_id,
			semester_id,
			start_date,
			end_date,
			created_at
		`)

	for _, g := range groupSemesters {
		builder = builder.Values(
			g.GroupID,
			g.SemesterID,
			g.StartDate,
			g.EndDate,
			g.CreatedAt,
		)
	}

	return builder.
		ReturningID().
		QueryxContext(r.db, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&groupSemesters[*i].ID); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r *GroupRepository) GetGroupSemesterByID(ctx context.Context, id uint64) (*domain.GroupSemester, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			id,
			group_id,
			semester_id,
			start_date,
			end_date,
			created_at,
			modified_at
		`).
		From("education_group_semester").
		WhereID(id).
		IsDeleted(false).
		ToSQL()
	if err != nil {
		return nil, err
	}

	groupSemester := &domain.GroupSemester{}
	if err = sqlx.GetContext(ctx, r.db, groupSemester, query, args...); err != nil {
		return nil, err
	}

	return groupSemester, nil
}

func (r *GroupRepository) FetchGroupSemesters(ctx context.Context, filters f.Filters) (domain.GroupSemesters, error) {
	query, args, err := q.Select(ctx).
		Columns(`
			id,
			group_id,
			semester_id,
			start_date,
			end_date,
			created_at,
			modified_at
		`).
		From("education_group_semester").
		ApplyQueryFilters(filters, q.Wheres{
			{
				Operator: q.And,
				Conditions: q.Conditions{
					{Name: domain.IDsParam, Column: "id", Operator: q.Equals},
					{Name: domain.GroupIDParam, Column: "group_id", Operator: q.Equals},
					{Name: domain.SemesterIDParam, Column: "semester_id", Operator: q.Equals},
				},
			},
		}).
		IsDeleted(false).
		OrderBy("id").
		ToSQL()
	if err != nil {
		return nil, err
	}

	groupSemesters := make(domain.GroupSemesters, 0)
	if err = sqlx.SelectContext(ctx, r.db, &groupSemesters, query, args...); err != nil {
		return nil, err
	}

	return groupSemesters, nil
}

func (r *GroupRepository) UpdateGroupSemesters(ctx context.Context, groupSemesters domain.GroupSemesters) error {
	builder := q.BatchUpdate(ctx).
		Table("education_group_semester").
		PrimaryKey("id").
		NumberOfRows(len(groupSemesters)).
		AddBigintColumn("id").
		AddBigintColumn("group_id").
		AddBigintColumn("semester_id").
		AddTimeStampColumn("start_date").
		AddTimeStampColumn("end_date").
		AddTimeStampColumn("modified_at")

	args := make([]any, 0, 6*len(groupSemesters))
	for _, g := range groupSemesters {
		args = append(args,
			g.ID,
			g.GroupID,
			g.SemesterID,
			g.StartDate,
			g.EndDate,
			g.ModifiedAt,
		)
	}

	return builder.
		ReturningCreatedAt().
		QueryxContext(r.db, args, func(rows *sqlx.Rows, i *int) error {
			for ; rows.Next(); *i++ {
				if err := rows.Scan(&groupSemesters[*i].CreatedAt); err != nil {
					return err
				}
			}
			return nil
		})
}

func (r *GroupRepository) DeleteGroupSemesters(ctx context.Context, ids []uint64) error {
	query, args, err := q.SoftDelete(ctx).
		Table("education_group_semester").
		WhereID(ids).
		ToSQL()
	if err != nil {
		return err
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
