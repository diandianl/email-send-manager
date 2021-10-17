<template>
  <div>
    <BasicLayout>
      <template #wrapper>
        <el-row :gutter="30" class="mb10">
          <el-col>
            <el-card shadow="always">
              <el-button type="text" style="padding-left: 15px;" @click="handleSendBatch">
                <el-card class="box-card" style="background-color: #9fe27e;align-items: center;" shadow="always" align="center">
                  <i class="el-icon-s-promotion" style="font-size: 32px;margin-right: 10px;" />
                  <span style="font-size: 24px;">发邮件</span>
                </el-card>
              </el-button>
              <el-button type="text" @click="handleOpenSetting">
                <el-card class="box-card" style="background-color: #ced2da;align-items: center;" shadow="always" align="center">
                  <i class="el-icon-setting" style="font-size: 32px;margin-right: 10px;" />
                  <span style="font-size: 24px;">邮件发送配置</span>
                </el-card>
              </el-button>
            </el-card>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="16">
            <el-card class="box-card">
              <div slot="header" class="clearfix">
                <span>发送记录</span>
              </div>
              <el-card class="box-card">
                <el-form ref="queryForm" :model="queryParams" :inline="true" label-width="100px">
                  <el-form-item label="模板">
                    <el-select v-model="queryParams.template_id" placeholder="按模板筛选" clearable size="small">
                      <el-option
                        v-for="tpl in templates"
                        :key="tpl.id"
                        :label="tpl.name"
                        :value="tpl.id"
                      />
                    </el-select>
                  </el-form-item>
                  <el-form-item label="客户邮箱">
                    <el-input
                      v-model="queryParams.customer_email"
                      placeholder="请输入客户邮箱"
                      clearable
                      size="small"
                      @keyup.enter.native="handleQuery"
                    />
                  </el-form-item>
                  <el-form-item label="发送结果状态">
                    <el-select v-model="queryParams.status" placeholder="发送状态" clearable size="small">
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

                <el-table v-loading="loading" :data="records" border>
                  <el-table-column label="模板" align="center">
                    <template slot-scope="scope">
                      {{ scope.row.template.name }}
                    </template>
                  </el-table-column>
                  <el-table-column label="客户" align="center">
                    <template slot-scope="scope">
                      {{ scope.row.customer.email }}
                    </template>
                  </el-table-column>
                  <el-table-column label="状态" align="center">
                    <template slot-scope="scope">
                      <el-popover
                        v-if="scope.row.status === 2"
                        placement="top-start"
                        title="失败原因"
                        width="200"
                        trigger="hover"
                        :content="scope.row.reason"
                      >
                        <el-tag slot="reference" type="danger" disable-transitions>失败</el-tag>
                      </el-popover>
                      <el-tag v-if="scope.row.status === 1" type="success" disable-transitions>成功</el-tag>
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
                        v-if="scope.row.status === 2"
                        size="mini"
                        type="text"
                        icon="el-icon-s-promotion"
                        @click="handleReSend(scope.row)"
                      >重发</el-button>
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
              </el-card>
            </el-card>
          </el-col>
          <el-col :span="8">
            <el-card v-if="current" class="box-card" :body-style="{ padding: '0' }">
              <div slot="header" class="clearfix">
                <span>当前任务</span>
                <el-button style="float: right" type="primary" plain size="small" @click="handleCancelSend">结束当前任务</el-button>
              </div>
              <div>
                <p>开始于: {{ current.start_at }}</p>
                <div class="tags">
                  <el-tag size="small">总计({{ current.total }})</el-tag>
                  <el-tag size="small" type="success">成功({{ current.success }})</el-tag>
                  <el-tag size="small" type="danger">失败({{ current.failure }})</el-tag>
                </div>
                <p><el-tag v-if="current.error" type="danger">任务失败: {{ current.error }}</el-tag></p>
              </div>
              <el-progress :percentage="current.percentage" />
            </el-card>
          </el-col>
        </el-row>

        <el-dialog title="邮件发送配置" :visible.sync="openSetting">
          <el-form ref="setting" :model="cfg" :rules="rules" label-width="100px">
            <el-form-item label="邮件服务器" prop="host">
              <el-input v-model="cfg.host" placeholder="请输入邮件服务器地址" />
            </el-form-item>
            <el-form-item label="端口" prop="port">
              <el-input-number v-model="cfg.port" controls-position="right" placeholder="请输入监听端口" />
            </el-form-item>
            <el-form-item label="用户名" prop="username">
              <el-input v-model="cfg.username" placeholder="请输入用户名" />
            </el-form-item>
            <el-form-item label="密码" prop="password">
              <el-input v-model="cfg.password" placeholder="请输入密码" show-password />
            </el-form-item>
            <el-form-item label="发送间隔(ms)" prop="interval">
              <el-input-number v-model="cfg.interval" :step="100" controls-position="right" placeholder="请输入发送间隔时间" />
            </el-form-item>
          </el-form>
          <div slot="footer" class="dialog-footer">
            <el-button type="primary" @click="submitForm">确 定</el-button>
            <el-button @click="cancel">取 消</el-button>
          </div>
        </el-dialog>
        <el-dialog title="发送邮件" :visible.sync="openSending">
          <el-form ref="sending" :model="sendBatch" :rules="sendBatchRules" label-width="100px">
            <el-form-item label="使用模板" prop="template_id">
              <el-select v-model="sendBatch.template_id" placeholder="请选择模板" clearable size="small">
                <el-option
                  v-for="tpl in templates"
                  :key="tpl.id"
                  :label="tpl.name"
                  :value="tpl.id"
                />
              </el-select>
            </el-form-item>
            <el-form-item label="客户选择模式">
              <el-switch
                v-model="sendBatch.include"
                :active-value="true"
                active-text="包含"
                :inactive-value="false"
                inactive-text="排除"
              />
            </el-form-item>
            <el-form-item label="选择客户" prop="customers">
              <el-select
                v-model="sendBatch.customer_ids"
                multiple
                filterable
                remote
                reserve-keyword
                placeholder="请输入关键词"
                :remote-method="findCustomers"
                :loading="loading"
              >
                <el-option
                  v-for="c in customers"
                  :key="c.id"
                  :label="c.email"
                  :value="c.id"
                />
              </el-select>
            </el-form-item>
          </el-form>
          <div slot="footer" class="dialog-footer">
            <el-button type="primary" @click="doSend">开始发送</el-button>
            <el-button @click="cancelSend">取 消</el-button>
          </div>
        </el-dialog>
      </template>
    </BasicLayout>
  </div>
</template>

<script>
import { querySettings, upsertSettings } from '@/api/setting'
import { queryRecord, delRecord } from '@/api/record'
import { queryTemplate } from '@/api/template'
import { queryCustomer } from '@/api/customer'
import { queryCurrent, doSendBatch, cancelSend } from '@/api/send-batch'

export default {
  name: 'SendManage',
  components: {
  },
  data() {
    return {
      queryParams: {
        pageIndex: 1,
        pageSize: 10,
        template_id: undefined,
        customer_email: undefined,
        status: undefined
      },
      loading: false,
      total: 0,
      records: [],
      templates: [],
      customers: [],
      // 总条数
      openSetting: false,
      cfg: {},
      statusOptions: [{ label: '成功', value: 1 }, { label: '失败', value: 2 }],
      // 表单校验
      rules: {
        host: [
          { required: true, message: '邮件服务地址不能为空', trigger: 'blur' }
        ],
        port: [
          { required: true, message: '邮件服务端口不能为空', trigger: 'blur' }
        ],
        username: [
          { required: true, message: '邮件服务登录账号不能为空', trigger: 'blur' }
        ],
        password: [
          { required: true, message: '邮件服务登录密码不能为空', trigger: 'blur' }
        ]
      },
      openSending: false,
      sendBatch: {
        template_id: undefined,
        include: undefined,
        customer_ids: []
      },
      sendBatchRules: {
        template_id: [
          { required: true, message: '请选择模板', trigger: 'blur' }
        ]
      },
      current: undefined
    }
  },
  created() {
    this.getList()
    this.loadTemplates()
    this.loadCurrentTask()
  },
  methods: {
    getList() {
      this.loading = true
      queryRecord(this.queryParams).then(response => {
        this.records = response.data.list
        this.total = response.data.pagination.total
        this.loading = false
      })
    },
    loadTemplates() {
      queryTemplate({ pageSize: -1, lite: true }).then(response => {
        this.templates = response.data.list
      }).catch(err => {
        this.msgError(err)
      })
    },
    loadCurrentTask() {
      queryCurrent().then(response => {
        const cur = response.data
        if (cur) {
          if (cur.total > 0) {
            cur.percentage = ((cur.success + cur.failure) / cur.total) * 100
          } else {
            cur.percentage = 0
          }
          if (!cur.error) {
            setTimeout(() => {
              this.loadCurrentTask()
            }, 1000)
          }
          const last = this.current && this.current.success || 0
          if ((cur.success - last) > 5) {
            this.getList()
          }
        }
        this.current = cur
      }).catch(err => {
        this.msgError(err)
      })
    },
    loadSettings() {
      return querySettings('email_send_setting').then(response => {
        const email_send_settings = response.data && response.data.value || {}
        this.cfg = email_send_settings
        return email_send_settings
      }).catch(err => {
        this.msgError(err)
      })
    },
    findCustomers(keyword) {
      if (keyword !== '') {
        this.loading = true
        queryCustomer({ keyword, pageSize: 10 }).then(response => {
          this.customers = response.data.list
          this.loading = false
        }).catch(err => {
          this.msgError(err)
        })
      } else {
        this.customers = []
      }
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
    handleOpenSetting() {
      this.loadSettings().then(() => {
        this.openSetting = true
      })
    },
    submitForm() {
      this.$refs['setting'].validate(valid => {
        if (valid) {
          upsertSettings({ key: 'email_send_setting', value: this.cfg }).then(response => {
            this.openSetting = false
          }).catch(err => {
            this.msgError(err)
          })
        }
      })
    },
    cancel() {
      this.openSetting = false
      this.reset()
    },
    // 表单重置
    reset() {
      this.cfg = {
        host: undefined,
        port: undefined,
        username: undefined,
        password: undefined,
        interval: undefined
      }
      this.resetForm('setting')
    },
    handleSendBatch() {
      this.loadSettings().then(settings => {
        if (settings && settings.host) {
          this.openSending = true
        } else {
          this.$confirm('还未配置邮件服务器，是否现在开始配置?', '警告', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }).then(() => {
            this.handleOpenSetting()
          }).catch(() => {})
        }
      })
    },
    doSend() {
      this.$refs['sending'].validate(valid => {
        if (valid) {
          doSendBatch(this.sendBatch).then(response => {
            this.openSending = false
            this.loadCurrentTask()
          }).catch(err => {
            this.msgError(err)
          })
        }
      })
    },
    cancelSend() {
      this.openSending = false
      this.resetSend()
    },
    // 表单重置
    resetSend() {
      this.sendBatch = {
        template_id: undefined,
        include: undefined,
        customer_ids: []
      }
      this.resetForm('sending')
    },
    /** 删除按钮操作 */
    handleDelete(row) {
      this.$confirm('操作不可恢复，是否确认删除?', '警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(function() {
        delRecord(row.id).then((response) => {
          this.open = false
          this.getList()
        }).catch(err => {
          this.msgError(err)
        })
      }).catch(() => {})
    },
    handleReSend(row) {
      this.sendBatch = {
        template_id: row.template_id,
        include: true,
        customer_ids: [row.customer_id]
      }
      const cs = this.customers || []
      if (!cs.find(c => c.id === row.customer_id)) {
        cs.push({ id: row.customer.id, email: row.customer.email })
      }
      this.handleSendBatch()
    },
    handleCancelSend() {
      if (!this.current.error) {
        this.$confirm('当前发送任务还未完成，是否确定取消?', '警告', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }).then(() => {
          return cancelSend()
        }).then(response => {
          this.loadCurrentTask()
        }).catch(err => {
          this.msgError(err)
        })
      } else {
        cancelSend().then(response => {
          this.loadCurrentTask()
        }).catch(err => {
          this.msgError(err)
        })
      }
    }
  }
}
</script>

<style lang="scss" scoped>
  .clearfix:before,
  .clearfix:after {
    display: table;
    content: "";
  }
  .clearfix:after {
    clear: both
  }
  .tags .el-tag {
    margin-right: 5px;
  }
</style>
