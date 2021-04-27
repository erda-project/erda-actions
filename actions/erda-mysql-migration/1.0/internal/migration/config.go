// Copyright (c) 2021 Terminus, Inc.
//
// This program is free software: you can use, redistribute, and/or modify
// it under the terms of the GNU Affero General Public License, version 3
// or later ("AGPL"), as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package migration

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/erda-project/erda/pkg/envconf"
	"github.com/erda-project/erda/pkg/sqllint/linters"
	"github.com/erda-project/erda/pkg/sqllint/rules"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	addonListURI       = "/api/addons?type=project&value=%s"
	addonDetailURI     = "/api/addons/%s"
	addonReferencesURI = "/api/addons/%s/actions/references"
)

var configuration *Conf

type Conf struct {
	// basic envs
	OrgID             uint64 `env:"DICE_ORG_ID" required:"true"`
	CiOpenapiToken    string `env:"DICE_OPENAPI_TOKEN" required:"true"`
	DiceOpenapiPrefix string `env:"DICE_OPENAPI_ADDR" required:"true"`
	ProjectName       string `env:"DICE_PROJECT_NAME" required:"true"`
	AppName           string `env:"DICE_APPLICATION_NAME" required:"true"`
	ProjectID         int64  `env:"DICE_PROJECT_ID" required:"true"`
	AppID             uint64 `env:"DICE_APPLICATION_ID" required:"true"`
	Workspace         string `env:"DICE_WORKSPACE" required:"true"`

	PipelineDebugMode bool `env:"PIPELINE_DEBUG_MODE"`

	// action envs
	WorkDir       string   `env:"ACTION_WORKDIR"`
	Database_     string   `env:"ACTION_DATABASE"`
	MigrationDir_ string   `env:"ACTION_MIGRATIONDIR"`
	NeedMySQLLint bool     `env:"ACTION_NEEDMYSQLLINT"`
	Modules_      []string `env:"ACTION_MODULES"`
	Linters       struct {
		// Rules
		BooleanFieldLinter          bool `json:"boolean_field_linter"`
		CharsetLinter               bool `json:"charset_linter"`
		ColumnNameLinter            bool `json:"column_name_linter"`
		ColumnCommentLinter         bool `json:"column_comment_linter"`
		DDLDMLLinter                bool `json:"ddldml_linter"`
		DestructLinter              bool `json:"destruct_linter"`
		FloatDoubleLinter           bool `json:"float_double_linter"`
		ForeignKeyLinter            bool `json:"foreign_key_linter"`
		IndexLengthLinter           bool `json:"index_length_linter"`
		IndexNameLinter             bool `json:"index_name_linter"`
		KeywordsLinter              bool `json:"keywords_linter"`
		IDExistsLinter              bool `json:"id_exists_linter"`
		IDTypeLinter                bool `json:"id_type_linter"`
		IDIsPrimaryLinter           bool `json:"id_is_primary_linter"`
		CreatedAtExistsLinter       bool `json:"created_at_exists_linter"`
		CreatedAtTypeLinter         bool `json:"created_at_type_linter"`
		CreatedAtDefaultValueLinter bool `json:"created_at_default_value_linter"`
		UpdatedAtExistsLinter       bool `json:"updated_at_exists_linter"`
		UpdatedAtTypeLinter         bool `json:"updated_at_type_linter"`
		UpdatedAtDefaultValueLinter bool `json:"updated_at_default_value_linter"`
		UpdatedAtOnUpdateLinter     bool `json:"updated_at_on_update_linter"`
		NotNullLinter               bool `json:"not_null_linter"`
		TableCommentLinter          bool `json:"table_comment_linter"`
		TableNameLinter             bool `json:"table_name_linter"`
		VarcharLengthLinter         bool `json:"varchar_length_linter"`

		// allowed ddl stmt
		CreateDatabaseStmt                   bool `json:"create_database_stmt"`
		AlterDatabaseStmt                    bool `json:"alter_database_stmt"`
		DropDatabaseStmt                     bool `json:"drop_database_stmt"`
		CreateTableStmt                      bool `json:"create_table_stmt"`
		DropTableStmt                        bool `json:"drop_table_stmt"`
		DropSequenceStmt                     bool `json:"drop_sequence_stmt"`
		RenameTableStmt                      bool `json:"rename_table_stmt"`
		CreateViewStmt                       bool `json:"create_view_stmt"`
		CreateSequenceStmt                   bool `json:"create_sequence_stmt"`
		CreateIndexStmt                      bool `json:"create_index_stmt"`
		DropIndexStmt                        bool `json:"drop_index_stmt"`
		LockTablesStmt                       bool `json:"lock_tables_stmt"`
		UnlockTablesStmt                     bool `json:"unlock_tables_stmt"`
		CleanupTableLockStmt                 bool `json:"cleanup_table_lock_stmt"`
		RepairTableStmt                      bool `json:"repair_table_stmt"`
		TruncateTableStmt                    bool `json:"truncate_table_stmt"`
		RecoverTableStmt                     bool `json:"recover_table_stmt"`
		FlashBackTableStmt                   bool `json:"flash_back_table_stmt"`
		AlterTableOption                     bool `json:"alter_table_option"`
		AlterTableAddColumns                 bool `json:"alter_table_add_columns"`
		AlterTableAddConstraint              bool `json:"alter_table_add_constraint"`
		AlterTableDropColumn                 bool `json:"alter_table_drop_column"`
		AlterTableDropPrimaryKey             bool `json:"alter_table_drop_primary_key"`
		AlterTableDropIndex                  bool `json:"alter_table_drop_index"`
		AlterTableDropForeignKey             bool `json:"alter_table_drop_foreign_key"`
		AlterTableModifyColumn               bool `json:"alter_table_modify_column"`
		AlterTableChangeColumn               bool `json:"alter_table_change_column"`
		AlterTableRenameColumn               bool `json:"alter_table_rename_column"`
		AlterTableRenameTable                bool `json:"alter_table_rename_table"`
		AlterTableAlterColumn                bool `json:"alter_table_alter_column"`
		AlterTableLock                       bool `json:"alter_table_lock"`
		AlterTableAlgorithm                  bool `json:"alter_table_algorithm"`
		AlterTableRenameIndex                bool `json:"alter_table_rename_index"`
		AlterTableForce                      bool `json:"alter_table_force"`
		AlterTableAddPartitions              bool `json:"alter_table_add_partitions"`
		AlterTableCoalescePartitions         bool `json:"alter_table_coalesce_partitions"`
		AlterTableDropPartition              bool `json:"alter_table_drop_partition"`
		AlterTableTruncatePartition          bool `json:"alter_table_truncate_partition"`
		AlterTablePartition                  bool `json:"alter_table_partition"`
		AlterTableEnableKeys                 bool `json:"alter_table_enable_keys"`
		AlterTableDisableKeys                bool `json:"alter_table_disable_keys"`
		AlterTableRemovePartitioning         bool `json:"alter_table_remove_partitioning"`
		AlterTableWithValidation             bool `json:"alter_table_with_validation"`
		AlterTableWithoutValidation          bool `json:"alter_table_without_validation"`
		AlterTableSecondaryLoad              bool `json:"alter_table_secondary_load"`
		AlterTableSecondaryUnload            bool `json:"alter_table_secondary_unload"`
		AlterTableRebuildPartition           bool `json:"alter_table_rebuild_partition"`
		AlterTableReorganizePartition        bool `json:"alter_table_reorganize_partition"`
		AlterTableCheckPartitions            bool `json:"alter_table_check_partitions"`
		AlterTableExchangePartition          bool `json:"alter_table_exchange_partition"`
		AlterTableOptimizePartition          bool `json:"alter_table_optimize_partition"`
		AlterTableRepairPartition            bool `json:"alter_table_repair_partition"`
		AlterTableImportPartitionTablespace  bool `json:"alter_table_import_partition_tablespace"`
		AlterTableDiscardPartitionTablespace bool `json:"alter_table_discard_partition_tablespace"`
		AlterTableAlterCheck                 bool `json:"alter_table_alter_check"`
		AlterTableDropCheck                  bool `json:"alter_table_drop_check"`
		AlterTableImportTablespace           bool `json:"alter_table_import_tablespace"`
		AlterTableDiscardTablespace          bool `json:"alter_table_discard_tablespace"`
		AlterTableIndexInvisible             bool `json:"alter_table_index_invisible"`
		AlterTableOrderByColumns             bool `json:"alter_table_order_by_columns"`
		AlterTableSetTiFlashReplica          bool `json:"alter_table_set_ti_flash_replica"`

		// allowed dml stmt
		SelectStmt      bool `json:"select_stmt"`
		UnionStmt       bool `json:"union_stmt"`
		LoadDataStmt    bool `json:"load_data_stmt"`
		InsertStmt      bool `json:"insert_stmt"`
		DeleteStmt      bool `json:"delete_stmt"`
		UpdateStmt      bool `json:"update_stmt"`
		ShowStmt        bool `json:"show_stmt"`
		SplitRegionStmt bool `json:"split_region_stmt"`
	} `env:"ACTION_LINTERS"`

	MetaFilename_ string `env:"METAFILE"`

	dsn string
}

func Configuration() *Conf {
	if configuration == nil {
		configuration = new(Conf)
		if err := envconf.Load(configuration); err != nil {
			logrus.Fatalf("failed to load configuration, err: %v", err)
		}
		if configuration.PipelineDebugMode {
			logrus.SetLevel(logrus.DebugLevel)
		} else {
			logrus.SetLevel(logrus.InfoLevel)
		}

		if err := configuration.getDSN(); err != nil {
			logrus.Fatalf("failed to get MySQL addon DSN, err: %v", err)
		}
	}

	return configuration
}

// DSN gets MySQL DSN
func (c *Conf) DSN() string {
	return c.dsn
}

// SandboxDSN gets sandbox DSN
func (c *Conf) SandboxDSN() string {
	return "root:12345678@(localhost:3306)/"
}

// MigrationDir gets migration scripts direction like .dice/migrations or migrations
func (c *Conf) MigrationDir() string {
	return c.MigrationDir_
}

// AppVersion gets application version
func (c *Conf) AppVersion() string {
	return ""
}

// BaseVersion gets base version
func (c *Conf) BaseVersion() string {
	return ""
}

// DebugSQL gets weather to debug SQL executing
func (c *Conf) DebugSQL() bool {
	return c.PipelineDebugMode
}

func (c *Conf) Database() string {
	return c.Database_
}

func (c *Conf) Workdir() string {
	return c.WorkDir
}

func (c *Conf) MetaFilename() string {
	return c.MetaFilename_
}

func (c *Conf) NeedErdaMySQLLint() bool {
	return c.NeedMySQLLint
}

func (c *Conf) Modules() []string {
	var modules []string
	for _, v := range c.Modules_ {
		ss := strings.Split(v, ",")
		for _, vv := range ss {
			modules = append(modules, strings.TrimSpace(vv))
		}
	}
	return modules
}

func (c *Conf) Rules() []rules.Ruler {
	var rulers []rules.Ruler
	var ddls, dmls uint64
	var ddlsM = map[uint64]bool{
		linters.CreateDatabaseStmt:                   c.Linters.CreateDatabaseStmt,
		linters.AlterDatabaseStmt:                    c.Linters.AlterDatabaseStmt,
		linters.DropDatabaseStmt:                     c.Linters.DropDatabaseStmt,
		linters.CreateTableStmt:                      c.Linters.CreateTableStmt,
		linters.DropTableStmt:                        c.Linters.DropTableStmt,
		linters.DropSequenceStmt:                     c.Linters.DropSequenceStmt,
		linters.RenameTableStmt:                      c.Linters.RenameTableStmt,
		linters.CreateViewStmt:                       c.Linters.CreateViewStmt,
		linters.CreateSequenceStmt:                   c.Linters.CreateSequenceStmt,
		linters.CreateIndexStmt:                      c.Linters.CreateIndexStmt,
		linters.DropIndexStmt:                        c.Linters.DropIndexStmt,
		linters.LockTablesStmt:                       c.Linters.LockTablesStmt,
		linters.UnlockTablesStmt:                     c.Linters.UnlockTablesStmt,
		linters.CleanupTableLockStmt:                 c.Linters.CleanupTableLockStmt,
		linters.RepairTableStmt:                      c.Linters.RepairTableStmt,
		linters.TruncateTableStmt:                    c.Linters.TruncateTableStmt,
		linters.RecoverTableStmt:                     c.Linters.RecoverTableStmt,
		linters.FlashBackTableStmt:                   c.Linters.FlashBackTableStmt,
		linters.AlterTableOption:                     c.Linters.AlterTableOption,
		linters.AlterTableAddColumns:                 c.Linters.AlterTableAddColumns,
		linters.AlterTableAddConstraint:              c.Linters.AlterTableAddConstraint,
		linters.AlterTableDropColumn:                 c.Linters.AlterTableDropColumn,
		linters.AlterTableDropPrimaryKey:             c.Linters.AlterTableDropPrimaryKey,
		linters.AlterTableDropIndex:                  c.Linters.AlterTableDropIndex,
		linters.AlterTableDropForeignKey:             c.Linters.AlterTableDropForeignKey,
		linters.AlterTableModifyColumn:               c.Linters.AlterTableModifyColumn,
		linters.AlterTableChangeColumn:               c.Linters.AlterTableChangeColumn,
		linters.AlterTableRenameColumn:               c.Linters.AlterTableRenameColumn,
		linters.AlterTableRenameTable:                c.Linters.AlterTableRenameTable,
		linters.AlterTableAlterColumn:                c.Linters.AlterTableAlterColumn,
		linters.AlterTableLock:                       c.Linters.AlterTableLock,
		linters.AlterTableAlgorithm:                  c.Linters.AlterTableAlgorithm,
		linters.AlterTableRenameIndex:                c.Linters.AlterTableRenameIndex,
		linters.AlterTableForce:                      c.Linters.AlterTableForce,
		linters.AlterTableAddPartitions:              c.Linters.AlterTableAddPartitions,
		linters.AlterTableCoalescePartitions:         c.Linters.AlterTableCoalescePartitions,
		linters.AlterTableDropPartition:              c.Linters.AlterTableDropPartition,
		linters.AlterTableTruncatePartition:          c.Linters.AlterTableTruncatePartition,
		linters.AlterTablePartition:                  c.Linters.AlterTablePartition,
		linters.AlterTableEnableKeys:                 c.Linters.AlterTableEnableKeys,
		linters.AlterTableDisableKeys:                c.Linters.AlterTableDisableKeys,
		linters.AlterTableRemovePartitioning:         c.Linters.AlterTableRemovePartitioning,
		linters.AlterTableWithValidation:             c.Linters.AlterTableWithValidation,
		linters.AlterTableWithoutValidation:          c.Linters.AlterTableWithoutValidation,
		linters.AlterTableSecondaryLoad:              c.Linters.AlterTableSecondaryLoad,
		linters.AlterTableSecondaryUnload:            c.Linters.AlterTableSecondaryUnload,
		linters.AlterTableRebuildPartition:           c.Linters.AlterTableRebuildPartition,
		linters.AlterTableReorganizePartition:        c.Linters.AlterTableReorganizePartition,
		linters.AlterTableCheckPartitions:            c.Linters.AlterTableCheckPartitions,
		linters.AlterTableExchangePartition:          c.Linters.AlterTableExchangePartition,
		linters.AlterTableOptimizePartition:          c.Linters.AlterTableOptimizePartition,
		linters.AlterTableRepairPartition:            c.Linters.AlterTableRepairPartition,
		linters.AlterTableImportPartitionTablespace:  c.Linters.AlterTableImportPartitionTablespace,
		linters.AlterTableDiscardPartitionTablespace: c.Linters.AlterTableDiscardPartitionTablespace,
		linters.AlterTableAlterCheck:                 c.Linters.AlterTableAlterCheck,
		linters.AlterTableDropCheck:                  c.Linters.AlterTableDropCheck,
		linters.AlterTableImportTablespace:           c.Linters.AlterTableImportTablespace,
		linters.AlterTableDiscardTablespace:          c.Linters.AlterTableDiscardTablespace,
		linters.AlterTableIndexInvisible:             c.Linters.AlterTableIndexInvisible,
		linters.AlterTableOrderByColumns:             c.Linters.AlterTableOrderByColumns,
		linters.AlterTableSetTiFlashReplica:          c.Linters.AlterTableSetTiFlashReplica,
	}
	for k, v := range ddlsM {
		if v {
			ddls |= k
		}
	}

	var dmlsM = map[uint64]bool{
		linters.SelectStmt:      c.Linters.SelectStmt,
		linters.UnionStmt:       c.Linters.UnionStmt,
		linters.LoadDataStmt:    c.Linters.LoadDataStmt,
		linters.InsertStmt:      c.Linters.InsertStmt,
		linters.DeleteStmt:      c.Linters.DeleteStmt,
		linters.UpdateStmt:      c.Linters.UpdateStmt,
		linters.ShowStmt:        c.Linters.ShowStmt,
		linters.SplitRegionStmt: c.Linters.SplitRegionStmt,
	}
	for k, v := range dmlsM {
		if v {
			dmls |= k
		}
	}

	var lintersConfigM = map[string]bool{
		"BooleanFieldLinter":          c.Linters.BooleanFieldLinter,
		"CharsetLinter":               c.Linters.CharsetLinter,
		"ColumnNameLinter":            c.Linters.ColumnNameLinter,
		"ColumnCommentLinter":         c.Linters.ColumnCommentLinter,
		"FloatDoubleLinter":           c.Linters.FloatDoubleLinter,
		"ForeignKeyLinter":            c.Linters.ForeignKeyLinter,
		"IndexLengthLinter":           c.Linters.IndexLengthLinter,
		"IndexNameLinter":             c.Linters.IndexNameLinter,
		"KeywordsLinter":              c.Linters.KeywordsLinter,
		"IDExistsLinter":              c.Linters.IDExistsLinter,
		"IDTypeLinter":                c.Linters.IDTypeLinter,
		"IDIsPrimaryLinter":           c.Linters.IDIsPrimaryLinter,
		"CreatedAtExistsLinter":       c.Linters.CreatedAtExistsLinter,
		"CreatedAtTypeLinter":         c.Linters.CreatedAtTypeLinter,
		"CreatedAtDefaultValueLinter": c.Linters.CreatedAtDefaultValueLinter,
		"UpdatedAtExistsLinter":       c.Linters.UpdatedAtExistsLinter,
		"UpdatedAtTypeLinter":         c.Linters.UpdatedAtTypeLinter,
		"UpdatedAtDefaultValueLinter": c.Linters.UpdatedAtDefaultValueLinter,
		"UpdatedAtOnUpdateLinter":     c.Linters.UpdatedAtOnUpdateLinter,
		"NotNullLinter":               c.Linters.NotNullLinter,
		"TableCommentLinter":          c.Linters.TableCommentLinter,
		"TableNameLinter":             c.Linters.TableNameLinter,
		"VarcharLengthLinter":         c.Linters.VarcharLengthLinter,
	}
	var lintersM = map[string]rules.Ruler{
		"BooleanFieldLinter":          linters.NewBooleanFieldLinter,
		"CharsetLinter":               linters.NewCharsetLinter,
		"ColumnNameLinter":            linters.NewColumnNameLinter,
		"ColumnCommentLinter":         linters.NewColumnCommentLinter,
		"FloatDoubleLinter":           linters.NewFloatDoubleLinter,
		"ForeignKeyLinter":            linters.NewForeignKeyLinter,
		"IndexLengthLinter":           linters.NewIndexLengthLinter,
		"IndexNameLinter":             linters.NewIndexNameLinter,
		"KeywordsLinter":              linters.NewKeywordsLinter,
		"IDExistsLinter":              linters.NewIDExistsLinter,
		"IDTypeLinter":                linters.NewIDTypeLinter,
		"IDIsPrimaryLinter":           linters.NewIDIsPrimaryLinter,
		"CreatedAtExistsLinter":       linters.NewCreatedAtExistsLinter,
		"CreatedAtTypeLinter":         linters.NewCreatedAtDefaultValueLinter,
		"CreatedAtDefaultValueLinter": linters.NewCreatedAtDefaultValueLinter,
		"UpdatedAtExistsLinter":       linters.NewUpdatedAtExistsLinter,
		"UpdatedAtTypeLinter":         linters.NewUpdatedAtTypeLinter,
		"UpdatedAtDefaultValueLinter": linters.NewUpdatedAtDefaultValueLinter,
		"UpdatedAtOnUpdateLinter":     linters.NewUpdatedAtOnUpdateLinter,
		"NotNullLinter":               linters.NewNotNullLinter,
		"TableCommentLinter":          linters.NewTableCommentLinter,
		"TableNameLinter":             linters.NewTableNameLinter,
		"VarcharLengthLinter":         linters.NewVarcharLengthLinter,
	}

	for name, ok := range lintersConfigM {
		if ok {
			rulers = append(rulers, lintersM[name])
		}
	}

	rulers = append(rulers, linters.NewAllowedStmtLinter(ddls, dmls))

	return rulers
}

func (c *Conf) getDSN() error {
	// 查找项目下所有的 addon 实例
	url := c.DiceOpenapiPrefix + fmt.Sprintf(addonListURI, strconv.FormatUint(uint64(c.ProjectID), 10))
	header := map[string][]string{"authorization": {c.CiOpenapiToken}}
	list, err := getAddonList(url, header)
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return errors.Errorf("there is no addon in the project, projectID: %v", c.ProjectID)
	}

	// filter mysql with the workspace
	var mysqlAddons []GetAddonsListResponseDataEle
	for _, addon := range list {
		if strings.EqualFold(addon.AddonName, "mysql") && strings.EqualFold(addon.Workspace, c.Workspace) {
			mysqlAddons = append(mysqlAddons, addon)
		}
	}
	if len(mysqlAddons) == 0 {
		return errors.Errorf("there is no MySQL addon on the current workspace %s", c.Workspace)
	}

	for _, addon := range mysqlAddons {
		url := c.DiceOpenapiPrefix + fmt.Sprintf(addonReferencesURI, addon.InstanceID)
		references, err := getAddonReferences(url, header)
		if err != nil {
			return err
		}

		for _, ref := range references {
			if ref.ApplicationID != c.AppID {
				continue
			}

			url := c.DiceOpenapiPrefix + fmt.Sprintf(addonDetailURI, addon.InstanceID)
			detail, err := getAddonDetail(url, header)
			if err != nil {
				return err
			}

			c.dsn = fmt.Sprintf("%s:%s@(%s:%s)/",
				detail.Config.MySQLUserName,
				detail.Config.MySQLPassword,
				detail.Config.MySQLHost,
				detail.Config.MySQLPort)

			return nil
		}
	}

	return errors.Errorf("mysql addon not found, applicationID: %v, workspace: %s", c.AppID, c.Workspace)
}
