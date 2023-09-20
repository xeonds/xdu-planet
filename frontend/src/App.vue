<template>
  <el-scrollbar style="width: 100vw; height: 100vh">
    <el-header>
      <el-menu mode="horizontal" :ellipsis="false" router>
        <el-menu-item class="title" route="/">XDU-Planet</el-menu-item>
        <el-menu-item index="1" route="/">主页</el-menu-item>
        <el-menu-item index="2" route="/member">成员</el-menu-item>
        <el-menu-item index="3" route="/about">关于</el-menu-item>
      </el-menu>
    </el-header>
    <el-main>
      <router-view />
    </el-main>
    <el-footer style="border-top: 1px solid #ccc; padding: 1rem">
      <el-row>
        <el-col :span="24">
          <el-link href="https://github.com/xeonds/xdu-planet">2023 © XDU-Planet</el-link>
          <el-text type="primary" style="margin: 1rem">|</el-text>
          <el-text>Powered by Vue.js</el-text>
        </el-col>
      </el-row>
    </el-footer>
  </el-scrollbar>
</template>

<style scoped>
.title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--el-color-primary);
}
.el-header {
  padding: 0;
}
.el-main {
  width: 100vw;
  max-width: 1000px;
  margin: 0 auto;
  min-height: calc(100vh - 64px - 64px);
}
</style>

<style>
:root {
  --el-color-primary: #ee7171;
  --el-menu-hover-bg-color: #ffe6e6;
}
html.dark {
  /* 自定义深色背景颜色 */
  --el-bg-color: #626aef;
}
.flex-grow {
  flex-grow: 1;
}
</style>
