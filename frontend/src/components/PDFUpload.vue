<template>
    <div class="pdf-upload">
      <el-upload
        class="upload-demo"
        drag
        action="/api/pdf/upload"
        :on-success="handleSuccess"
        :on-error="handleError"
        accept="application/pdf"
      >
        <i class="el-icon-upload"></i>
        <div class="el-upload__text">拖拽PDF文件或 <em>点击上传</em></div>
      </el-upload>
      
      <div class="search-section">
        <el-input
          v-model="searchQuery"
          placeholder="输入搜索内容"
        ></el-input>
        <el-button @click="handleSearch">搜索</el-button>
      </div>
      
      <div class="results-section">
        <el-table :data="searchResults">
          <el-table-column prop="text" label="内容"></el-table-column>
          <el-table-column prop="pageNumber" label="页码"></el-table-column>
          <el-table-column prop="similarity" label="相似度"></el-table-column>
        </el-table>
      </div>
    </div>
  </template>
  
  <script setup>
  import { ref } from 'vue'
  import axios from 'axios'
  
  const searchQuery = ref('')
  const searchResults = ref([])
  
  const handleSuccess = (response) => {
    ElMessage.success('PDF上传成功')
  }
  
  const handleError = (error) => {
    ElMessage.error('上传失败：' + error.message)
  }
  
  const handleSearch = async () => {
    try {
      const response = await axios.post('/api/vectors/search', {
        text: searchQuery.value,
        limit: 10
      })
      searchResults.value = response.data
    } catch (error) {
      ElMessage.error('搜索失败：' + error.message)
    }
  }
  </script>