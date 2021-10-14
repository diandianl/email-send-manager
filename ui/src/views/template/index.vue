<template>
  <BasicLayout>
    <template #wrapper>
      <el-card class="box-card">
        <el-form ref="queryForm" :model="queryParams" :inline="true" label-width="68px">
          <el-form-item label="名称" prop="email">
            <el-input
              v-model="queryParams.name"
              placeholder="请输入模板名称"
              clearable
              size="small"
              @keyup.enter.native="handleQuery"
            />
          </el-form-item>
          <el-form-item label="状态" prop="status">
            <el-select v-model="queryParams.status" placeholder="岗位状态" clearable size="small">
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
        </el-row>

        <el-table v-loading="loading" :data="list" border>
          <el-table-column label="名称" align="center" prop="name" />
          <el-table-column label="主题" align="center" prop="subject" />
          <el-table-column label="发件人邮箱" align="center" prop="from" />
          <el-table-column label="发件人名称" align="center" prop="from_name" />
          <el-table-column label="回复邮箱" align="center" prop="reply_to" />
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
                @click="handlePreview(scope.row)"
              >预览</el-button>
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
        <el-dialog :title="title" :visible.sync="open">
          <el-form ref="form" :model="form" :rules="rules" label-width="80px">
            <el-row :gutter="5">
              <el-col :span="10">
                <el-form-item label="模板名称" label-width="85px" prop="name">
                  <el-input v-model="form.name" placeholder="请输入模板名称" />
                </el-form-item>
              </el-col>
              <el-col :span="10" :offset="2">
                <el-form-item label="邮件主题" label-width="85px" prop="subject">
                  <el-input v-model="form.subject" placeholder="请输入邮件主题" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="5">
              <el-col :span="10">
                <el-form-item label="发件人邮箱" label-width="85px" prop="from">
                  <el-input v-model="form.from" placeholder="请输入发件人邮箱" />
                </el-form-item>
              </el-col>
              <el-col :span="10" :offset="2">
                <el-form-item label="发件人名称" label-width="85px" prop="from_name">
                  <el-input v-model="form.from_name" placeholder="请输入发件人名称" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="5">
              <el-col :span="10">
                <el-form-item label="回复邮箱" label-width="85px" prop="reply_to">
                  <el-input v-model="form.reply_to" placeholder="请输入回复邮箱" />
                </el-form-item>
              </el-col>
              <el-col :span="10" :offset="2">
                <el-form-item label="状态" prop="status" label-width="85px">
                  <el-switch v-model="form.status" :active-value="1" active-text="启用" :inactive-value="0" inactive-text="禁用" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="5" style="height: 20rem">
              <el-col>
                <el-divider content-position="left">邮件正文模板</el-divider>
                <quill-editor
                  v-model="form.content"
                  class="editor"
                  style="height: 100%"
                  :options="editorOption"
                />
              </el-col>
            </el-row>
          </el-form>
          <div slot="footer" class="dialog-footer">
            <el-button type="primary" @click="submitForm">确 定</el-button>
            <el-button @click="cancel">取 消</el-button>
          </div>
        </el-dialog>
        <el-dialog title="预览邮件正文" :visible.sync="preview">
          <div class="ql-editor" v-html="previewContent" />
        </el-dialog>
      </el-card>
    </template>
  </BasicLayout>
</template>

<script>
import { queryTemplate, getTemplate, updateTemplate, addTemplate, delTemplate } from '@/api/template'

import Quill from 'quill'
import 'quill/dist/quill.core.css'
import 'quill/dist/quill.snow.css'
import 'quill/dist/quill.bubble.css'
import { ImageDrop } from 'quill-image-drop-module'
import ImageResize from 'quill-image-resize-module'
Quill.register('modules/imageDrop', ImageDrop)
Quill.register('modules/imageResize', ImageResize)

import { quillEditor } from 'vue-quill-editor'

export default {
  name: 'TemplateManage',
  components: {
    quillEditor
  },
  data() {
    return {
      editorOption: {
        modules: {
          toolbar: [
            ['bold', 'italic', 'underline', 'strike'],
            ['blockquote', 'code-block'],
            [{ 'header': 1 }, { 'header': 2 }],
            [{ 'list': 'ordered' }, { 'list': 'bullet' }],
            [{ 'script': 'sub' }, { 'script': 'super' }],
            [{ 'indent': '-1' }, { 'indent': '+1' }],
            [{ 'direction': 'rtl' }],
            [{ 'size': ['small', false, 'large', 'huge'] }],
            [{ 'header': [1, 2, 3, 4, 5, 6, false] }],
            [{ 'font': [] }],
            [{ 'color': [] }, { 'background': [] }],
            [{ 'align': [] }],
            ['clean'],
            ['link', 'image', 'emoji']
          ],
          history: {
            delay: 1000,
            maxStack: 50,
            userOnly: false
          },
          imageDrop: true,
          imageResize: {
            displayStyles: {
              backgroundColor: 'black',
              border: 'none',
              color: 'white'
            },
            modules: ['Resize', 'DisplaySize', 'Toolbar']
          }
        },
        placeholder: '输入内容........'
      },
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
      // 预览模板正文
      preview: false,
      previewContent: undefined,
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
      form: {},
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
      queryTemplate(this.queryParams).then(response => {
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
        subject: undefined,
        from: undefined,
        from_name: undefined,
        reply_to: undefined,
        content: undefined,
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
      this.title = '新建邮件模板'
    },
    /** 预览操作 */
    handlePreview(row) {
      this.previewContent = row.content
      this.preview = true
    },
    /** 修改按钮操作 */
    handleUpdate(row) {
      this.reset()
      getTemplate(row.id).then(response => {
        this.form = response.data
        this.open = true
        this.title = '修改邮件模板'
      })
    },
    /** 提交按钮 */
    submitForm: function() {
      this.$refs['form'].validate(valid => {
        if (valid) {
          this.form.status = parseInt(this.form.status)
          if (this.form.id !== undefined) {
            updateTemplate(this.form, this.form.id).then(response => {
              this.open = false
              this.getList()
            }).catch(err => {
              this.msgError(err)
            })
          } else {
            addTemplate(this.form).then(response => {
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
        return delTemplate(row.id)
      }).then((response) => {
        this.open = false
        this.getList()
      }).catch(err => {
        this.msgError(err)
      })
    }
  }
}
</script>
