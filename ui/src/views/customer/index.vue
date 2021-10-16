<template>
  <BasicLayout>
    <template #wrapper>
      <el-card class="box-card">
        <el-form ref="queryForm" :model="queryParams" :inline="true" label-width="68px">
          <el-form-item label="名称" prop="email">
            <el-input
              v-model="queryParams.name"
              placeholder="请输入客户名称"
              clearable
              size="small"
              @keyup.enter.native="handleQuery"
            />
          </el-form-item>
          <el-form-item label="邮箱" prop="email">
            <el-input
              v-model="queryParams.email"
              placeholder="请输入客户邮箱"
              clearable
              size="small"
              @keyup.enter.native="handleQuery"
            />
          </el-form-item>
          <el-form-item label="状态" prop="status">
            <el-select v-model="queryParams.status" placeholder="客户状态" clearable size="small">
              <el-option
                v-for="o in statusOptions"
                :key="o.value"
                :label="o.label"
                :value="o.value"
              />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" icon="el-icon-search" size="mini" @click="handleQuery">搜索</el-button>
            <el-button icon="el-icon-refresh" size="mini" @click="resetQuery">重置</el-button>
          </el-form-item>
        </el-form>

        <el-row :gutter="10" class="mb8">
          <el-col :span="1.5">
            <el-button type="primary" icon="el-icon-plus" size="mini" @click="handleAdd">新增</el-button>
          </el-col>
          <el-col :span="1.5">
            <el-button type="warning" icon="el-icon-upload" size="mini" @click="openImport = true">导入</el-button>
          </el-col>
        </el-row>

        <el-table v-loading="loading" :data="list" border>
          <el-table-column label="名称" align="center" prop="name" />
          <el-table-column label="邮箱" align="center" prop="email" />
          <el-table-column label="状态" align="center" prop="status">
            <template slot-scope="scope">
              <el-tag
                :type="scope.row.status === 0 ? 'danger' : 'success'"
                disable-transitions
              >{{ scope.row.status === 0 ? '禁用' : '启用' }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="创建时间" align="center" prop="created_at" width="180">
            <template slot-scope="scope">
              <span>{{ parseTime(scope.row.created_at) }}</span>
            </template>
          </el-table-column>
          <el-table-column label="操作" align="center" class-name="small-padding fixed-width">
            <template slot-scope="scope">
              <el-button
                size="mini"
                type="text"
                icon="el-icon-edit"
                @click="handleUpdate(scope.row)"
              >修改</el-button>
              <el-button
                size="mini"
                type="text"
                icon="el-icon-delete"
                @click="handleDelete(scope.row)"
              >删除</el-button>
            </template>
          </el-table-column>
        </el-table>

        <pagination
          v-show="total>0"
          :total="total"
          :page.sync="queryParams.pageIndex"
          :limit.sync="queryParams.pageSize"
          @pagination="getList"
        />

        <!-- 添加或修改岗位对话框 -->
        <el-dialog :title="title" :visible.sync="open" width="500px">
          <el-form ref="form" :model="form" :rules="rules" label-width="80px">
            <el-form-item label="客户名称" prop="name">
              <el-input v-model="form.name" placeholder="请输入客户名称" />
            </el-form-item>
            <el-form-item label="客户邮箱" prop="email">
              <el-input v-model="form.email" placeholder="请输入客户邮箱" />
            </el-form-item>
            <el-form-item label="状态" prop="status">
              <el-switch v-model="form.status" :active-value="1" active-text="启用" :inactive-value="0" inactive-text="禁用" />
            </el-form-item>
          </el-form>
          <div slot="footer" class="dialog-footer">
            <el-button type="primary" @click="submitForm">确 定</el-button>
            <el-button @click="cancel">取 消</el-button>
          </div>
        </el-dialog>
        <el-dialog title="批量导入客户信息" :visible.sync="openImport">
          <el-upload
            class="upload-demo"
            drag
            accept=".xls,.xlsx"
            :action="importUrl"
            :on-success="handleImportSuccess"
            :on-error="handleImportFailure"
          >
            <i class="el-icon-upload" />
            <div class="el-upload__text">将文件拖到此处，或<em>点击上传</em></div>
            <div slot="tip" class="el-upload__tip">只能上传 Excel 文件，且不超过500kb</div>
          </el-upload>
          <div slot="footer" class="dialog-footer">
            <el-button @click="openImport = false">关 闭</el-button>
          </div>
        </el-dialog>
      </el-card>
    </template>
  </BasicLayout>
</template>

<script>
import { queryCustomer, getCustomer, updateCustomer, addCustomer, delCustomer, importUrl } from '@/api/customer'

export default {
  name: 'CustomerManage',
  data() {
    return {
      importUrl: importUrl,
      // 遮罩层
      loading: true,
      // 选中数组
      ids: [],
      // 非单个禁用
      single: true,
      // 非多个禁用
      multiple: true,
      // 总条数
      total: 0,
      // 岗位表格数据
      list: [],
      // 弹出层标题
      title: '',
      // 是否显示弹出层
      open: false,
      // 是否显示导入对话框
      openImport: false,
      // 状态数据字典
      statusOptions: [{ label: '启用', value: 1 }, { label: '禁用', value: 0 }],
      // 查询参数
      queryParams: {
        pageIndex: 1,
        pageSize: 10,
        name: undefined,
        email: undefined,
        status: undefined
      },
      // 表单参数
      form: {
        id: undefined,
        name: undefined,
        email: undefined,
        status: 1
      },
      // 表单校验
      rules: {
        name: [
          { required: true, message: '客户名称不能为空', trigger: 'blur' }
        ],
        email: [
          { required: true, message: '客户邮箱不能为空', trigger: 'blur' }
        ]
      }
    }
  },
  created() {
    this.getList()
  },
  methods: {
    /** 查询岗位列表 */
    getList() {
      this.loading = true
      queryCustomer(this.queryParams).then(response => {
        this.list = response.data.list
        this.total = response.data.pagination.total
        this.loading = false
      })
    },
    // 岗位状态字典翻译
    statusFormat(row) {
      return this.selectDictLabel(this.statusOptions, row.status)
    },
    // 取消按钮
    cancel() {
      this.open = false
      this.reset()
    },
    // 表单重置
    reset() {
      this.form = {
        id: undefined,
        name: undefined,
        email: undefined,
        status: 1
      }
      this.resetForm('form')
    },
    /** 搜索按钮操作 */
    handleQuery() {
      this.queryParams.pageIndex = 1
      this.getList()
    },
    /** 重置按钮操作 */
    resetQuery() {
      this.resetForm('queryForm')
      this.handleQuery()
    },
    // 多选框选中数据
    handleSelectionChange(selection) {
      this.ids = selection.map(item => item.postId)
      this.single = selection.length !== 1
      this.multiple = !selection.length
    },
    /** 新增按钮操作 */
    handleAdd() {
      this.reset()
      this.open = true
      this.title = '添加客户信息'
    },
    /** 修改按钮操作 */
    handleUpdate(row) {
      this.reset()
      getCustomer(row.id).then(response => {
        this.form = response.data
        this.open = true
        this.title = '修改客户信息'
      })
    },
    /** 提交按钮 */
    submitForm: function() {
      this.$refs['form'].validate(valid => {
        if (valid) {
          if (this.form.id !== undefined) {
            updateCustomer(this.form, this.form.id).then(response => {
              this.open = false
              this.getList()
            }).catch(err => {
              this.msgError(err)
            })
          } else {
            addCustomer(this.form).then(response => {
              this.open = false
              this.getList()
            }).catch(err => {
              this.msgError(err)
            })
          }
        }
      })
    },
    /** 删除按钮操作 */
    handleDelete(row) {
      this.$confirm('操作不可恢复，是否确认删除?', '警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(function() {
        return delCustomer(row.id)
      }).then((response) => {
        this.open = false
        this.getList()
      }).catch(err => {
        this.msgError(err)
      })
    },
    handleImportSuccess(response) {
      this.msgInfo('导入成功')
      this.getList()
    },
    handleImportFailure(err) {
      this.msgError(err)
    }
  }
}
</script>
