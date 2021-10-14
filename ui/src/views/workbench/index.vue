<template>
  <div>
    <BasicLayout>
      <template #wrapper>
        <el-row :gutter="30" class="mb10">
          <el-col :span="5">
            <el-card class="box-card" shadow="always" align="center">
              <i class="el-icon-position" style="font-size: 32px;margin-right: 10px;"></i>
              <span style="font-size: 24px;">发邮件</span>
            </el-card>
          </el-col>
          <el-col :span="5">
            <el-card class="box-card" shadow="always" align="center">
              <el-button type="success" plain icon="el-icon-setting" @click="handleOpenSetting">邮件发送配置</el-button>
            </el-card>
          </el-col>
        </el-row>

        <el-card class="box-card">
          <div slot="header" class="clearfix">
            <span>发送记录</span>
          </div>
          <div class="monitor">
          </div>
        </el-card>

        <el-dialog title="邮件发送配置" :visible.sync="openSetting">
          <el-form label-width="100px">
            <el-form-item label="邮件服务器" prop="host">
              <el-input v-model="cfg.host" placeholder="请输入邮件服务器地址" />
            </el-form-item>
            <el-form-item label="端口" prop="port">
              <el-input v-model="cfg.port" placeholder="请输入监听端口" />
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
      </template>
    </BasicLayout>
  </div>
</template>

<script>
export default {
  name: 'SendManage',
  components: {
  },
  data() {
    return {
      openSetting: false,
      cfg: {}
    }
  },
  created() {
    this.getServerInfo()
  },
  methods: {
    getServerInfo() {
    },
    handleOpenSetting() {
      console.error('fuck')
      this.openSetting = true
    },
    submitForm() {},
    cancel() {}
  }
}
</script>

<style lang="scss" scoped>
  .el-card__body{
    padding: 20px 20px 0 20px!important;
  }
</style>
