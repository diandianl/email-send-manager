<template>
  <BasicLayout>
    <template #wrapper>
      <el-card class="box-card">
        <el-tabs tab-position="left" style="height: 100%;">
          <el-tab-pane label="系统配置">
            <el-form label-width="80px">
              <el-form-item label="监听地址" prop="host">
                <el-input v-model="sys.host" placeholder="请输入监听地址" />
              </el-form-item>
              <el-form-item label="监听端口" prop="port">
                <el-input v-model="sys.port" placeholder="请输入监听端口" />
              </el-form-item>
            </el-form>
            <div slot="footer" class="dialog-footer">
              <el-button type="primary" @click="submitForm">确 定</el-button>
              <el-button @click="cancel">取 消</el-button>
            </div>
          </el-tab-pane>
          <el-tab-pane label="邮件发送配置">
          </el-tab-pane>
        </el-tabs>
      </el-card>
    </template>
  </BasicLayout>
</template>

<script>
export default {
  name: 'SysConfigSet',
  components: {
  },
  data() {
    return {
      // 遮罩层
      loading: true,
      sys: {
        host: '0.0.0.0',
        port: 10099,
      },
      // 参数
      cfg: {}
    }
  },
  created() {
    this.loadConfig()
  },
  methods: {
    loadConfig() {
      this.loading = true
      getSetConfig().then(response => {
        this.configList = response.data
        this.loading = false
        this.fillFormData(this.formConf, this.configList)
        // 更新表单
        this.key2 = +new Date()
      }
      )
    },
    setUrl(url) {
      const data = {
        sys_app_logo: ''
      }
      data.sys_app_logo = url
      // 回填数据
      this.fillFormData(this.formConf, data)
      // 更新表单
      this.key2 = +new Date()
    },
    // 参数系统内置字典翻译
    typeFormat(row, column) {
      return this.selectDictLabel(this.typeOptions, row.configType)
    },
    /** 搜索按钮操作 */
    handleQuery() {
      this.queryParams.pageIndex = 1
      this.getList()
    },
    fillFormData(form, data) {
      form.fields.forEach(item => {
        const val = data[item.__vModel__]
        if (val) {
          item.__config__.defaultValue = val
        }
      })
    },
    bind(key, data) {
      this.setUrl(data)
    },
    sumbitForm2(data) {
      var list = []
      var i = 0
      for (var key in data) {
        list[i] = { 'configKey': key, 'configValue': data[key] }
        i++
      }
      updateSetConfig(list).then(response => {
        if (response.code === 200) {
          this.msgSuccess(response.msg)
          this.open = false
          this.getList()
        } else {
          this.msgError(response.msg)
        }
      }
      )
    }
  }
}
</script>
