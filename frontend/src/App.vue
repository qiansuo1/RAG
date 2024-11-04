<template>
    <el-container>
      <el-header>
        <el-menu mode="horizontal" router>
          <el-menu-item index="/">首页</el-menu-item>
        </el-menu>
      </el-header>
      
      <el-main>
        <router-view></router-view>
      </el-main>
    </el-container>
  </template>
  
  <script setup>
  </script>