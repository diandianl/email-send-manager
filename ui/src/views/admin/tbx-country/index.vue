
<template>
  <BasicLayout>
    <template #wrapper>
      <el-card class="box-card">
        <el-form ref="queryForm" :model="queryParams" :inline="true" label-width="68px">
          <el-form-item label="编码" prop="code">
            <el-input v-model="queryParams.code" placeholder="请输入编码" clearable size="small" @keyup.enter.native="handleQuery" />
          </el-form-item>
          <el-form-item label="中文名称" prop="cName">
            <el-input v-model="queryParams.cName" placeholder="请输入中文名称" clearable size="small" @keyup.enter.native="handleQuery" />
          </el-form-item>
          <el-form-item label="English Name" prop="eName">
            <el-input v-model="queryParams.eName" placeholder="请输入English Name" clearable size="small" @keyup.enter.native="handleQuery" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" icon="el-icon-search" size="mini" @click="handleQuery">搜索</el-button>
            <el-button icon="el-icon-refresh" size="mini" @click="resetQuery">重置</el-button>
          </el-form-item>
        </el-form>

        <el-row :gutter="10" class="mb8">
          <el-col :span="1.5">
            <el-button v-permisaction="['admin:tbxCountry:add']" type="primary" icon="el-icon-plus" size="mini" @click="handleAdd">新增</el-button>
          </el-col>
          <el-col :span="1.5">
            <el-button v-permisaction="['admin:tbxCountry:edit']" type="success" icon="el-icon-edit" size="mini" :disabled="single" @click="handleUpdate">修改</el-button>
          </el-col>
          <el-col :span="1.5">
            <el-button v-permisaction="['admin:tbxCountry:remove']" type="danger" icon="el-icon-delete" size="mini" :disabled="multiple" @click="handleDelete">删除</el-button>
          </el-col>
        </el-row>

        <el-table v-loading="loading" :data="tbxCountryList" @selection-change="handleSelectionChange">
          <el-table-column type="selection" width="55" align="center" />
          <el-table-column label="编码" align="center" prop="code" :show-overflow-tooltip="true" />
          <el-table-column label="中文名称" align="center" prop="cName" :show-overflow-tooltip="true" />
          <el-table-column label="English Name" align="center" prop="eName" :show-overflow-tooltip="true" />
          <el-table-column label="English Name 2" align="center" prop="eName2" :show-overflow-tooltip="true" />
          <el-table-column label="编码2" align="center" prop="code2" :show-overflow-tooltip="true" />

          <el-table-column label="操作" align="center" class-name="small-padding fixed-width">
            <template slot-scope="scope">
              <el-button v-permisaction="['admin:tbxCountry:edit']" size="mini" type="text" icon="el-icon-edit" @click="handleUpdate(scope.row)">修改</el-button>
              <el-button v-permisaction="['admin:tbxCountry:remove']" size="mini" type="text" icon="el-icon-delete" @click="handleDelete(scope.row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>

        <pagination v-show="total>0" :total="total" :page.sync="queryParams.pageIndex" :limit.sync="queryParams.pageSize" @pagination="getList" />

        <!-- 添加或修改对话框 -->
        <el-dialog :title="title" :visible.sync="open" width="500px">
          <el-form ref="form" :model="form" :rules="rules" label-width="80px">
            <el-form-item label="编码" prop="code">
              <el-input v-model="form.code" placeholder="编码" :disabled="isEdit" />
            </el-form-item>
            <el-form-item label="中文名称" prop="cName">
              <el-input v-model="form.cName" placeholder="中文名称" :disabled="isEdit" />
            </el-form-item>
            <el-form-item label="English Name" prop="eName">
              <el-input v-model="form.eName" placeholder="English Name" :disabled="isEdit" />
            </el-form-item>
            <el-form-item label="English Name 2" prop="eName2">
              <el-input v-model="form.eName2" placeholder="English Name 2" :disabled="isEdit" />
            </el-form-item>
            <el-form-item label="编码2" prop="code2">
              <el-input v-model="form.code2" placeholder="编码2" :disabled="isEdit" />
            </el-form-item>
            <el-form-item label="编码3" prop="code3">
              <el-input v-model="form.code3" placeholder="编码3" :disabled="isEdit" />
            </el-form-item>
          </el-form>
          <div slot="footer" class="dialog-footer">
            <el-button type="primary" @click="submitForm">确 定</el-button>
            <el-button @click="cancel">取 消</el-button>
          </div>
        </el-dialog>
      </el-card>
    </template>
  </BasicLayout>
</template>

<script>
import { addTbxCountry, delTbxCountry, getTbxCountry, listTbxCountry, updateTbxCountry } from '@/api/admin/tbx-country'

export default {
  name: 'TbxCountry',
  components: {
  },
  data() {
    return {
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
      // 弹出层标题
      title: '',
      // 是否显示弹出层
      open: false,
      isEdit: false,
      // 类型数据字典
      typeOptions: [],
      tbxCountryList: [],

      // 关系表类型

      // 查询参数
      queryParams: {
        pageIndex: 1,
        pageSize: 10,
        code: undefined,
        cName: undefined,
        eName: undefined

      },
      // 表单参数
      form: {
      },
      // 表单校验
      rules: {
        code: [{ required: true, message: '编码不能为空', trigger: 'blur' }],
        cName: [{ required: true, message: '中文名称不能为空', trigger: 'blur' }],
        eName: [{ required: true, message: 'English Name不能为空', trigger: 'blur' }]
      }
    }
  },
  created() {
    this.getList()
  },
  methods: {
    /** 查询参数列表 */
    getList() {
      this.loading = true
      listTbxCountry(this.addDateRange(this.queryParams, this.dateRange)).then(response => {
        this.tbxCountryList = response.data.list
        this.total = response.data.count
        this.loading = false
      })
    },
    // 取消按钮
    cancel() {
      this.open = false
      this.reset()
    },
    // 表单重置
    reset() {
      this.form = {

        code: undefined,
        cName: undefined,
        eName: undefined,
        eName2: undefined,
        code2: undefined,
        code3: undefined
      }
      this.resetForm('form')
    },
    getImgList: function() {
      this.form[this.fileIndex] = this.$refs['fileChoose'].resultList[0].fullUrl
    },
    fileClose: function() {
      this.fileOpen = false
    },
    // 关系
    // 文件
    /** 搜索按钮操作 */
    handleQuery() {
      this.queryParams.pageIndex = 1
      this.getList()
    },
    /** 重置按钮操作 */
    resetQuery() {
      this.dateRange = []
      this.resetForm('queryForm')
      this.handleQuery()
    },
    /** 新增按钮操作 */
    handleAdd() {
      this.reset()
      this.open = true
      this.title = '添加国家编码'
      this.isEdit = false
    },
    // 多选框选中数据
    handleSelectionChange(selection) {
      this.ids = selection.map(item => item.code)
      this.single = selection.length !== 1
      this.multiple = !selection.length
    },
    /** 修改按钮操作 */
    handleUpdate(row) {
      this.reset()
      const code = row.code || this.ids
      getTbxCountry(code).then(response => {
        this.form = response.data
        this.open = true
        this.title = '修改国家编码'
        this.isEdit = true
      })
    },
    /** 提交按钮 */
    submitForm: function() {
      this.$refs['form'].validate(valid => {
        if (valid) {
          if (this.isEdit) {
            updateTbxCountry(this.form).then(response => {
              if (response.code === 200) {
                this.msgSuccess(response.msg)
                this.open = false
                this.getList()
              } else {
                this.msgError(response.msg)
              }
            })
          } else {
            addTbxCountry(this.form).then(response => {
              if (response.code === 200) {
                this.msgSuccess(response.msg)
                this.open = false
                this.getList()
              } else {
                this.msgError(response.msg)
              }
            })
          }
        }
      })
    },
    /** 删除按钮操作 */
    handleDelete(row) {
      var Ids = (row.code && [row.code]) || this.ids

      this.$confirm('是否确认删除编号为"' + Ids + '"的数据项?', '警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(function() {
        return delTbxCountry({ 'ids': Ids })
      }).then((response) => {
        if (response.code === 200) {
          this.msgSuccess(response.msg)
          this.open = false
          this.getList()
        } else {
          this.msgError(response.msg)
        }
      }).catch(function() {
      })
    }
  }
}
</script>
