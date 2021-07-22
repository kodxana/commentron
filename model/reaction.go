// Code generated by SQLBoiler (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package model

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/volatiletech/sqlboiler/queries/qmhelper"
	"github.com/volatiletech/sqlboiler/strmangle"
)

// Reaction is an object representing the database table.
type Reaction struct {
	ID             uint64      `boil:"id" json:"id" toml:"id" yaml:"id"`
	CommentID      string      `boil:"comment_id" json:"comment_id" toml:"comment_id" yaml:"comment_id"`
	ChannelID      null.String `boil:"channel_id" json:"channel_id,omitempty" toml:"channel_id" yaml:"channel_id,omitempty"`
	ClaimID        string      `boil:"claim_id" json:"claim_id" toml:"claim_id" yaml:"claim_id"`
	ReactionTypeID uint64      `boil:"reaction_type_id" json:"reaction_type_id" toml:"reaction_type_id" yaml:"reaction_type_id"`
	CreatedAt      time.Time   `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt      time.Time   `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`
	IsFlagged      bool        `boil:"is_flagged" json:"is_flagged" toml:"is_flagged" yaml:"is_flagged"`

	R *reactionR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L reactionL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ReactionColumns = struct {
	ID             string
	CommentID      string
	ChannelID      string
	ClaimID        string
	ReactionTypeID string
	CreatedAt      string
	UpdatedAt      string
	IsFlagged      string
}{
	ID:             "id",
	CommentID:      "comment_id",
	ChannelID:      "channel_id",
	ClaimID:        "claim_id",
	ReactionTypeID: "reaction_type_id",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
	IsFlagged:      "is_flagged",
}

// Generated where

var ReactionWhere = struct {
	ID             whereHelperuint64
	CommentID      whereHelperstring
	ChannelID      whereHelpernull_String
	ClaimID        whereHelperstring
	ReactionTypeID whereHelperuint64
	CreatedAt      whereHelpertime_Time
	UpdatedAt      whereHelpertime_Time
	IsFlagged      whereHelperbool
}{
	ID:             whereHelperuint64{field: "`reaction`.`id`"},
	CommentID:      whereHelperstring{field: "`reaction`.`comment_id`"},
	ChannelID:      whereHelpernull_String{field: "`reaction`.`channel_id`"},
	ClaimID:        whereHelperstring{field: "`reaction`.`claim_id`"},
	ReactionTypeID: whereHelperuint64{field: "`reaction`.`reaction_type_id`"},
	CreatedAt:      whereHelpertime_Time{field: "`reaction`.`created_at`"},
	UpdatedAt:      whereHelpertime_Time{field: "`reaction`.`updated_at`"},
	IsFlagged:      whereHelperbool{field: "`reaction`.`is_flagged`"},
}

// ReactionRels is where relationship names are stored.
var ReactionRels = struct {
	Channel      string
	Comment      string
	ReactionType string
}{
	Channel:      "Channel",
	Comment:      "Comment",
	ReactionType: "ReactionType",
}

// reactionR is where relationships are stored.
type reactionR struct {
	Channel      *Channel
	Comment      *Comment
	ReactionType *ReactionType
}

// NewStruct creates a new relationship struct
func (*reactionR) NewStruct() *reactionR {
	return &reactionR{}
}

// reactionL is where Load methods for each relationship are stored.
type reactionL struct{}

var (
	reactionAllColumns            = []string{"id", "comment_id", "channel_id", "claim_id", "reaction_type_id", "created_at", "updated_at", "is_flagged"}
	reactionColumnsWithoutDefault = []string{"comment_id", "channel_id", "claim_id", "reaction_type_id"}
	reactionColumnsWithDefault    = []string{"id", "created_at", "updated_at", "is_flagged"}
	reactionPrimaryKeyColumns     = []string{"id"}
)

type (
	// ReactionSlice is an alias for a slice of pointers to Reaction.
	// This should generally be used opposed to []Reaction.
	ReactionSlice []*Reaction

	reactionQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	reactionType                 = reflect.TypeOf(&Reaction{})
	reactionMapping              = queries.MakeStructMapping(reactionType)
	reactionPrimaryKeyMapping, _ = queries.BindMapping(reactionType, reactionMapping, reactionPrimaryKeyColumns)
	reactionInsertCacheMut       sync.RWMutex
	reactionInsertCache          = make(map[string]insertCache)
	reactionUpdateCacheMut       sync.RWMutex
	reactionUpdateCache          = make(map[string]updateCache)
	reactionUpsertCacheMut       sync.RWMutex
	reactionUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single reaction record from the query.
func (q reactionQuery) One(exec boil.Executor) (*Reaction, error) {
	o := &Reaction{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: failed to execute a one query for reaction")
	}

	return o, nil
}

// All returns all Reaction records from the query.
func (q reactionQuery) All(exec boil.Executor) (ReactionSlice, error) {
	var o []*Reaction

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model: failed to assign all query results to Reaction slice")
	}

	return o, nil
}

// Count returns the count of all Reaction records in the query.
func (q reactionQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to count reaction rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q reactionQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model: failed to check if reaction exists")
	}

	return count > 0, nil
}

// Channel pointed to by the foreign key.
func (o *Reaction) Channel(mods ...qm.QueryMod) channelQuery {
	queryMods := []qm.QueryMod{
		qm.Where("claim_id=?", o.ChannelID),
	}

	queryMods = append(queryMods, mods...)

	query := Channels(queryMods...)
	queries.SetFrom(query.Query, "`channel`")

	return query
}

// Comment pointed to by the foreign key.
func (o *Reaction) Comment(mods ...qm.QueryMod) commentQuery {
	queryMods := []qm.QueryMod{
		qm.Where("comment_id=?", o.CommentID),
	}

	queryMods = append(queryMods, mods...)

	query := Comments(queryMods...)
	queries.SetFrom(query.Query, "`comment`")

	return query
}

// ReactionType pointed to by the foreign key.
func (o *Reaction) ReactionType(mods ...qm.QueryMod) reactionTypeQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=?", o.ReactionTypeID),
	}

	queryMods = append(queryMods, mods...)

	query := ReactionTypes(queryMods...)
	queries.SetFrom(query.Query, "`reaction_type`")

	return query
}

// LoadChannel allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (reactionL) LoadChannel(e boil.Executor, singular bool, maybeReaction interface{}, mods queries.Applicator) error {
	var slice []*Reaction
	var object *Reaction

	if singular {
		object = maybeReaction.(*Reaction)
	} else {
		slice = *maybeReaction.(*[]*Reaction)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &reactionR{}
		}
		if !queries.IsNil(object.ChannelID) {
			args = append(args, object.ChannelID)
		}

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &reactionR{}
			}

			for _, a := range args {
				if queries.Equal(a, obj.ChannelID) {
					continue Outer
				}
			}

			if !queries.IsNil(obj.ChannelID) {
				args = append(args, obj.ChannelID)
			}

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(qm.From(`channel`), qm.WhereIn(`claim_id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Channel")
	}

	var resultSlice []*Channel
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Channel")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for channel")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for channel")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Channel = foreign
		if foreign.R == nil {
			foreign.R = &channelR{}
		}
		foreign.R.Reactions = append(foreign.R.Reactions, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if queries.Equal(local.ChannelID, foreign.ClaimID) {
				local.R.Channel = foreign
				if foreign.R == nil {
					foreign.R = &channelR{}
				}
				foreign.R.Reactions = append(foreign.R.Reactions, local)
				break
			}
		}
	}

	return nil
}

// LoadComment allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (reactionL) LoadComment(e boil.Executor, singular bool, maybeReaction interface{}, mods queries.Applicator) error {
	var slice []*Reaction
	var object *Reaction

	if singular {
		object = maybeReaction.(*Reaction)
	} else {
		slice = *maybeReaction.(*[]*Reaction)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &reactionR{}
		}
		args = append(args, object.CommentID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &reactionR{}
			}

			for _, a := range args {
				if a == obj.CommentID {
					continue Outer
				}
			}

			args = append(args, obj.CommentID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(qm.From(`comment`), qm.WhereIn(`comment_id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Comment")
	}

	var resultSlice []*Comment
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Comment")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for comment")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for comment")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Comment = foreign
		if foreign.R == nil {
			foreign.R = &commentR{}
		}
		foreign.R.Reactions = append(foreign.R.Reactions, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.CommentID == foreign.CommentID {
				local.R.Comment = foreign
				if foreign.R == nil {
					foreign.R = &commentR{}
				}
				foreign.R.Reactions = append(foreign.R.Reactions, local)
				break
			}
		}
	}

	return nil
}

// LoadReactionType allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (reactionL) LoadReactionType(e boil.Executor, singular bool, maybeReaction interface{}, mods queries.Applicator) error {
	var slice []*Reaction
	var object *Reaction

	if singular {
		object = maybeReaction.(*Reaction)
	} else {
		slice = *maybeReaction.(*[]*Reaction)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &reactionR{}
		}
		args = append(args, object.ReactionTypeID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &reactionR{}
			}

			for _, a := range args {
				if a == obj.ReactionTypeID {
					continue Outer
				}
			}

			args = append(args, obj.ReactionTypeID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(qm.From(`reaction_type`), qm.WhereIn(`id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load ReactionType")
	}

	var resultSlice []*ReactionType
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice ReactionType")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for reaction_type")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for reaction_type")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.ReactionType = foreign
		if foreign.R == nil {
			foreign.R = &reactionTypeR{}
		}
		foreign.R.Reactions = append(foreign.R.Reactions, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.ReactionTypeID == foreign.ID {
				local.R.ReactionType = foreign
				if foreign.R == nil {
					foreign.R = &reactionTypeR{}
				}
				foreign.R.Reactions = append(foreign.R.Reactions, local)
				break
			}
		}
	}

	return nil
}

// SetChannel of the reaction to the related item.
// Sets o.R.Channel to related.
// Adds o to related.R.Reactions.
func (o *Reaction) SetChannel(exec boil.Executor, insert bool, related *Channel) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE `reaction` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, []string{"channel_id"}),
		strmangle.WhereClause("`", "`", 0, reactionPrimaryKeyColumns),
	)
	values := []interface{}{related.ClaimID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	queries.Assign(&o.ChannelID, related.ClaimID)
	if o.R == nil {
		o.R = &reactionR{
			Channel: related,
		}
	} else {
		o.R.Channel = related
	}

	if related.R == nil {
		related.R = &channelR{
			Reactions: ReactionSlice{o},
		}
	} else {
		related.R.Reactions = append(related.R.Reactions, o)
	}

	return nil
}

// RemoveChannel relationship.
// Sets o.R.Channel to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *Reaction) RemoveChannel(exec boil.Executor, related *Channel) error {
	var err error

	queries.SetScanner(&o.ChannelID, nil)
	if err = o.Update(exec, boil.Whitelist("channel_id")); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.R.Channel = nil
	if related == nil || related.R == nil {
		return nil
	}

	for i, ri := range related.R.Reactions {
		if queries.Equal(o.ChannelID, ri.ChannelID) {
			continue
		}

		ln := len(related.R.Reactions)
		if ln > 1 && i < ln-1 {
			related.R.Reactions[i] = related.R.Reactions[ln-1]
		}
		related.R.Reactions = related.R.Reactions[:ln-1]
		break
	}
	return nil
}

// SetComment of the reaction to the related item.
// Sets o.R.Comment to related.
// Adds o to related.R.Reactions.
func (o *Reaction) SetComment(exec boil.Executor, insert bool, related *Comment) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE `reaction` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, []string{"comment_id"}),
		strmangle.WhereClause("`", "`", 0, reactionPrimaryKeyColumns),
	)
	values := []interface{}{related.CommentID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.CommentID = related.CommentID
	if o.R == nil {
		o.R = &reactionR{
			Comment: related,
		}
	} else {
		o.R.Comment = related
	}

	if related.R == nil {
		related.R = &commentR{
			Reactions: ReactionSlice{o},
		}
	} else {
		related.R.Reactions = append(related.R.Reactions, o)
	}

	return nil
}

// SetReactionType of the reaction to the related item.
// Sets o.R.ReactionType to related.
// Adds o to related.R.Reactions.
func (o *Reaction) SetReactionType(exec boil.Executor, insert bool, related *ReactionType) error {
	var err error
	if insert {
		if err = related.Insert(exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE `reaction` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, []string{"reaction_type_id"}),
		strmangle.WhereClause("`", "`", 0, reactionPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.ReactionTypeID = related.ID
	if o.R == nil {
		o.R = &reactionR{
			ReactionType: related,
		}
	} else {
		o.R.ReactionType = related
	}

	if related.R == nil {
		related.R = &reactionTypeR{
			Reactions: ReactionSlice{o},
		}
	} else {
		related.R.Reactions = append(related.R.Reactions, o)
	}

	return nil
}

// Reactions retrieves all the records using an executor.
func Reactions(mods ...qm.QueryMod) reactionQuery {
	mods = append(mods, qm.From("`reaction`"))
	return reactionQuery{NewQuery(mods...)}
}

// FindReaction retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindReaction(exec boil.Executor, iD uint64, selectCols ...string) (*Reaction, error) {
	reactionObj := &Reaction{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `reaction` where `id`=?", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, reactionObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: unable to select from reaction")
	}

	return reactionObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Reaction) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model: no reaction provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(reactionColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	reactionInsertCacheMut.RLock()
	cache, cached := reactionInsertCache[key]
	reactionInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			reactionAllColumns,
			reactionColumnsWithDefault,
			reactionColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(reactionType, reactionMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(reactionType, reactionMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `reaction` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `reaction` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `reaction` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, reactionPrimaryKeyColumns))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	result, err := exec.Exec(cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "model: unable to insert into reaction")
	}

	var lastID int64
	var identifierCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	lastID, err = result.LastInsertId()
	if err != nil {
		return ErrSyncFail
	}

	o.ID = uint64(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == reactionMapping["ID"] {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.ID,
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.retQuery)
		fmt.Fprintln(boil.DebugWriter, identifierCols...)
	}

	err = exec.QueryRow(cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "model: unable to populate default values for reaction")
	}

CacheNoHooks:
	if !cached {
		reactionInsertCacheMut.Lock()
		reactionInsertCache[key] = cache
		reactionInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the Reaction.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Reaction) Update(exec boil.Executor, columns boil.Columns) error {
	var err error
	key := makeCacheKey(columns, nil)
	reactionUpdateCacheMut.RLock()
	cache, cached := reactionUpdateCache[key]
	reactionUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			reactionAllColumns,
			reactionPrimaryKeyColumns,
		)

		if len(wl) == 0 {
			return errors.New("model: unable to update reaction, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `reaction` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, reactionPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(reactionType, reactionMapping, append(wl, reactionPrimaryKeyColumns...))
		if err != nil {
			return err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err = exec.Exec(cache.query, values...)
	if err != nil {
		return errors.Wrap(err, "model: unable to update reaction row")
	}

	if !cached {
		reactionUpdateCacheMut.Lock()
		reactionUpdateCache[key] = cache
		reactionUpdateCacheMut.Unlock()
	}

	return nil
}

// UpdateAll updates all rows with the specified column values.
func (q reactionQuery) UpdateAll(exec boil.Executor, cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec(exec)
	if err != nil {
		return errors.Wrap(err, "model: unable to update all for reaction")
	}

	return nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ReactionSlice) UpdateAll(exec boil.Executor, cols M) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	if len(cols) == 0 {
		return errors.New("model: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), reactionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `reaction` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, reactionPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "model: unable to update all in reaction slice")
	}

	return nil
}

var mySQLReactionUniqueColumns = []string{
	"id",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Reaction) Upsert(exec boil.Executor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("model: no reaction provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(reactionColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLReactionUniqueColumns, o)

	if len(nzUniques) == 0 {
		return errors.New("cannot upsert with a table that cannot conflict on a unique column")
	}

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzUniques {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	reactionUpsertCacheMut.RLock()
	cache, cached := reactionUpsertCache[key]
	reactionUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			reactionAllColumns,
			reactionColumnsWithDefault,
			reactionColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			reactionAllColumns,
			reactionPrimaryKeyColumns,
		)

		if len(update) == 0 {
			return errors.New("model: unable to upsert reaction, could not build update column list")
		}

		ret = strmangle.SetComplement(ret, nzUniques)
		cache.query = buildUpsertQueryMySQL(dialect, "reaction", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `reaction` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(reactionType, reactionMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(reactionType, reactionMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	result, err := exec.Exec(cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "model: unable to upsert for reaction")
	}

	var lastID int64
	var uniqueMap []uint64
	var nzUniqueCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	lastID, err = result.LastInsertId()
	if err != nil {
		return ErrSyncFail
	}

	o.ID = uint64(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == reactionMapping["id"] {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(reactionType, reactionMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "model: unable to retrieve unique values for reaction")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.retQuery)
		fmt.Fprintln(boil.DebugWriter, nzUniqueCols...)
	}

	err = exec.QueryRow(cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "model: unable to populate default values for reaction")
	}

CacheNoHooks:
	if !cached {
		reactionUpsertCacheMut.Lock()
		reactionUpsertCache[key] = cache
		reactionUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single Reaction record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Reaction) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("model: no Reaction provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), reactionPrimaryKeyMapping)
	sql := "DELETE FROM `reaction` WHERE `id`=?"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "model: unable to delete from reaction")
	}

	return nil
}

// DeleteAll deletes all matching rows.
func (q reactionQuery) DeleteAll(exec boil.Executor) error {
	if q.Query == nil {
		return errors.New("model: no reactionQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec(exec)
	if err != nil {
		return errors.Wrap(err, "model: unable to delete all from reaction")
	}

	return nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ReactionSlice) DeleteAll(exec boil.Executor) error {
	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), reactionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `reaction` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, reactionPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "model: unable to delete all from reaction slice")
	}

	return nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Reaction) Reload(exec boil.Executor) error {
	ret, err := FindReaction(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ReactionSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := ReactionSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), reactionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `reaction`.* FROM `reaction` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, reactionPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model: unable to reload all in ReactionSlice")
	}

	*o = slice

	return nil
}

// ReactionExists checks if the Reaction row exists.
func ReactionExists(exec boil.Executor, iD uint64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `reaction` where `id`=? limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}

	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model: unable to check if reaction exists")
	}

	return exists, nil
}
