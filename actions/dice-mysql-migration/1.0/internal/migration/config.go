package migration

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/erda-project/erda-actions/actions/dice-mysql-migration/1.0/internal/log"
	"github.com/erda-project/erda/pkg/envconf"
	"github.com/erda-project/erda/pkg/sqllint"
	"github.com/pkg/errors"
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
	WorkDir       string `env:"ACTION_WORKDIR"`
	Database_     string `env:"ACTION_DATABASE"`
	MigrationDir_ string `env:"ACTION_MIGRATIONDIR"`
	NeedMySQLLint bool   `env:"ACTION_NEEDMYSQLLINT"`

	Linters struct {
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
			log.Fatalf("failed to load configuration, err: %v", err)
		}

		if err := configuration.getDSN(); err != nil {
			log.Fatalf("failed to get MySQL addon DSN, err: %v", err)
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

func (c *Conf) Rules() []sqllint.NewRule {
	var rules []sqllint.NewRule
	var ddls, dmls uint64
	if c.Linters.CreateDatabaseStmt {
		ddls |= sqllint.CreateDatabaseStmt
	}
	if c.Linters.AlterDatabaseStmt {
		ddls |= sqllint.AlterDatabaseStmt
	}
	if c.Linters.DropDatabaseStmt {
		ddls |= sqllint.DropDatabaseStmt
	}
	if c.Linters.CreateTableStmt {
		ddls |= sqllint.CreateTableStmt
	}
	if c.Linters.DropTableStmt {
		ddls |= sqllint.DropTableStmt
	}
	if c.Linters.DropSequenceStmt {
		ddls |= sqllint.DropSequenceStmt
	}
	if c.Linters.RenameTableStmt {
		ddls |= sqllint.RenameTableStmt
	}
	if c.Linters.CreateViewStmt {
		ddls |= sqllint.CreateViewStmt
	}
	if c.Linters.CreateSequenceStmt {
		ddls |= sqllint.CreateSequenceStmt
	}
	if c.Linters.CreateIndexStmt {
		ddls |= sqllint.CreateIndexStmt
	}
	if c.Linters.DropIndexStmt {
		ddls |= sqllint.DropIndexStmt
	}
	if c.Linters.LockTablesStmt {
		ddls |= sqllint.LockTablesStmt
	}
	if c.Linters.UnlockTablesStmt {
		ddls |= sqllint.UnlockTablesStmt
	}
	if c.Linters.CleanupTableLockStmt {
		ddls |= sqllint.CleanupTableLockStmt
	}
	if c.Linters.RepairTableStmt {
		ddls |= sqllint.RepairTableStmt
	}
	if c.Linters.TruncateTableStmt {
		ddls |= sqllint.TruncateTableStmt
	}
	if c.Linters.RecoverTableStmt {
		ddls |= sqllint.RecoverTableStmt
	}
	if c.Linters.FlashBackTableStmt {
		ddls |= sqllint.FlashBackTableStmt
	}
	if c.Linters.AlterTableOption {
		ddls |= sqllint.AlterTableOption
	}
	if c.Linters.AlterTableAddColumns {
		ddls |= sqllint.AlterTableAddColumns
	}
	if c.Linters.AlterTableAddConstraint {
		ddls |= sqllint.AlterTableAddConstraint
	}
	if c.Linters.AlterTableDropColumn {
		ddls |= sqllint.AlterTableDropColumn
	}
	if c.Linters.AlterTableDropPrimaryKey {
		ddls |= sqllint.AlterTableDropPrimaryKey
	}
	if c.Linters.AlterTableDropIndex {
		ddls |= sqllint.AlterTableDropIndex
	}
	if c.Linters.AlterTableDropForeignKey {
		ddls |= sqllint.AlterTableDropForeignKey
	}
	if c.Linters.AlterTableModifyColumn {
		ddls |= sqllint.AlterTableModifyColumn
	}
	if c.Linters.AlterTableChangeColumn {
		ddls |= sqllint.AlterTableChangeColumn
	}
	if c.Linters.AlterTableRenameColumn {
		ddls |= sqllint.AlterTableRenameColumn
	}
	if c.Linters.AlterTableRenameTable {
		ddls |= sqllint.AlterTableRenameTable
	}
	if c.Linters.AlterTableAlterColumn {
		ddls |= sqllint.AlterTableAlterColumn
	}
	if c.Linters.AlterTableLock {
		ddls |= sqllint.AlterTableLock
	}
	if c.Linters.AlterTableAlgorithm {
		ddls |= sqllint.AlterTableAlgorithm
	}
	if c.Linters.AlterTableRenameIndex {
		ddls |= sqllint.AlterTableRenameIndex
	}
	if c.Linters.AlterTableForce {
		ddls |= sqllint.AlterTableForce
	}
	if c.Linters.AlterTableAddPartitions {
		ddls |= sqllint.AlterTableAddPartitions
	}
	if c.Linters.AlterTableCoalescePartitions {
		ddls |= sqllint.AlterTableCoalescePartitions
	}
	if c.Linters.AlterTableDropPartition {
		ddls |= sqllint.AlterTableDropPartition
	}
	if c.Linters.AlterTableTruncatePartition {
		ddls |= sqllint.AlterTableTruncatePartition
	}
	if c.Linters.AlterTablePartition {
		ddls |= sqllint.AlterTablePartition
	}
	if c.Linters.AlterTableEnableKeys {
		ddls |= sqllint.AlterTableEnableKeys
	}
	if c.Linters.AlterTableDisableKeys {
		ddls |= sqllint.AlterTableDisableKeys
	}
	if c.Linters.AlterTableRemovePartitioning {
		ddls |= sqllint.AlterTableRemovePartitioning
	}
	if c.Linters.AlterTableWithValidation {
		ddls |= sqllint.AlterTableWithValidation
	}
	if c.Linters.AlterTableWithoutValidation {
		ddls |= sqllint.AlterTableWithoutValidation
	}
	if c.Linters.AlterTableSecondaryLoad {
		ddls |= sqllint.AlterTableSecondaryLoad
	}
	if c.Linters.AlterTableSecondaryUnload {
		ddls |= sqllint.AlterTableSecondaryUnload
	}
	if c.Linters.AlterTableRebuildPartition {
		ddls |= sqllint.AlterTableRebuildPartition
	}
	if c.Linters.AlterTableReorganizePartition {
		ddls |= sqllint.AlterTableReorganizePartition
	}
	if c.Linters.AlterTableCheckPartitions {
		ddls |= sqllint.AlterTableCheckPartitions
	}
	if c.Linters.AlterTableExchangePartition {
		ddls |= sqllint.AlterTableExchangePartition
	}
	if c.Linters.AlterTableOptimizePartition {
		ddls |= sqllint.AlterTableOptimizePartition
	}
	if c.Linters.AlterTableRepairPartition {
		ddls |= sqllint.AlterTableRepairPartition
	}
	if c.Linters.AlterTableImportPartitionTablespace {
		ddls |= sqllint.AlterTableImportPartitionTablespace
	}
	if c.Linters.AlterTableDiscardPartitionTablespace {
		ddls |= sqllint.AlterTableDiscardPartitionTablespace
	}
	if c.Linters.AlterTableAlterCheck {
		ddls |= sqllint.AlterTableAlterCheck
	}
	if c.Linters.AlterTableDropCheck {
		ddls |= sqllint.AlterTableDropCheck
	}
	if c.Linters.AlterTableImportTablespace {
		ddls |= sqllint.AlterTableImportTablespace
	}
	if c.Linters.AlterTableDiscardTablespace {
		ddls |= sqllint.AlterTableDiscardTablespace
	}
	if c.Linters.AlterTableIndexInvisible {
		ddls |= sqllint.AlterTableIndexInvisible
	}
	if c.Linters.AlterTableOrderByColumns {
		ddls |= sqllint.AlterTableOrderByColumns
	}
	if c.Linters.AlterTableSetTiFlashReplica {
		ddls |= sqllint.AlterTableSetTiFlashReplica
	}

	if c.Linters.SelectStmt {
		dmls |= sqllint.SelectStmt
	}
	if c.Linters.UnionStmt {
		dmls |= sqllint.UnionStmt
	}
	if c.Linters.LoadDataStmt {
		dmls |= sqllint.LoadDataStmt
	}
	if c.Linters.InsertStmt {
		dmls |= sqllint.InsertStmt
	}
	if c.Linters.DeleteStmt {
		dmls |= sqllint.DeleteStmt
	}
	if c.Linters.UpdateStmt {
		dmls |= sqllint.UpdateStmt
	}
	if c.Linters.ShowStmt {
		dmls |= sqllint.ShowStmt
	}
	if c.Linters.SplitRegionStmt {
		dmls |= sqllint.SplitRegionStmt
	}

	if c.Linters.BooleanFieldLinter {
		rules = append(rules, sqllint.NewBooleanFieldLinter)
	}
	if c.Linters.CharsetLinter {
		rules = append(rules, sqllint.NewCharsetLinter)
	}
	if c.Linters.ColumnNameLinter {
		rules = append(rules, sqllint.NewColumnNameLinter)
	}
	if c.Linters.ColumnCommentLinter {
		rules = append(rules, sqllint.NewColumnCommentLinter)
	}
	// if c.DDLDMLLinter {
	// 	rules = append(rules, sqllint.NewDDLDMLLinter)
	// }
	// if c.DestructLinter {
	// 	rules = append(rules, sqllint.NewDestructLinter)
	// }
	if c.Linters.FloatDoubleLinter {
		rules = append(rules, sqllint.NewFloatDoubleLinter)
	}
	if c.Linters.ForeignKeyLinter {
		rules = append(rules, sqllint.NewForeignKeyLinter)
	}
	if c.Linters.IndexLengthLinter {
		rules = append(rules, sqllint.NewIndexLengthLinter)
	}
	if c.Linters.IndexNameLinter {
		rules = append(rules, sqllint.NewIndexNameLinter)
	}
	if c.Linters.KeywordsLinter {
		rules = append(rules, sqllint.NewKeywordsLinter)
	}
	if c.Linters.IDExistsLinter {
		rules = append(rules, sqllint.NewIDExistsLinter)
	}
	if c.Linters.IDTypeLinter {
		rules = append(rules, sqllint.NewIDTypeLinter)
	}
	if c.Linters.IDIsPrimaryLinter {
		rules = append(rules, sqllint.NewIDIsPrimaryLinter)
	}
	if c.Linters.CreatedAtExistsLinter {
		rules = append(rules, sqllint.NewCreatedAtExistsLinter)
	}
	if c.Linters.CreatedAtTypeLinter {
		rules = append(rules, sqllint.NewCreatedAtTypeLinter)
	}
	if c.Linters.CreatedAtDefaultValueLinter {
		rules = append(rules, sqllint.NewCreatedAtDefaultValueLinter)
	}
	if c.Linters.UpdatedAtExistsLinter {
		rules = append(rules, sqllint.NewUpdatedAtExistsLinter)
	}
	if c.Linters.UpdatedAtTypeLinter {
		rules = append(rules, sqllint.NewUpdatedAtTypeLinter)
	}
	if c.Linters.UpdatedAtDefaultValueLinter {
		rules = append(rules, sqllint.NewUpdatedAtDefaultValueLinter)
	}
	if c.Linters.UpdatedAtOnUpdateLinter {
		rules = append(rules, sqllint.NewUpdatedAtOnUpdateLinter)
	}
	if c.Linters.NotNullLinter {
		rules = append(rules, sqllint.NewNotNullLinter)
	}
	if c.Linters.TableCommentLinter {
		rules = append(rules, sqllint.NewTableCommentLinter)
	}
	if c.Linters.TableNameLinter {
		rules = append(rules, sqllint.NewTableNameLinter)
	}
	if c.Linters.VarcharLengthLinter {
		rules = append(rules, sqllint.NewVarcharLengthLinter)
	}

	rules = append(rules, sqllint.NewAllowedStmt(ddls, dmls))

	return rules
}

func (c *Conf) getDSN() error {
	// 查找项目下所有的 addon 实例
	url := c.DiceOpenapiPrefix + fmt.Sprintf(addonListURI, strconv.FormatUint(uint64(c.ProjectID), 10))
	header := map[string][]string{"authorization": {c.CiOpenapiToken}}
	list, err := getAddonList(url, header)
	if err != nil {
		return err
	}

	// 查找 workspace 与流水线中 workspace 一致的 mysql addon 的引用列表
	for _, addon := range list {
		if strings.ToLower(addon.AddonName) != "mysql" {
			continue
		}
		if strings.EqualFold(addon.Workspace, c.Workspace) {
			continue
		}

		url := c.DiceOpenapiPrefix + fmt.Sprintf(addonReferencesURI, addon.InstanceID)
		references, err := getAddonReferences(url, header)
		if err != nil {
			return err
		}

		// 如果本应用引用了这个 mysql addon, 则获取这个 mysql addon 的配置详情
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

	return errors.Errorf("mysql addon not found, applicationID: %v, workspce: %s", c.AppID, c.Workspace)
}
