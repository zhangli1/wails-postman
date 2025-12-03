

<main>
 <div class="sidebar">
     <div class="sidebar-header">
         <div class="sidebar-title">
             <i class="fas fa-layer-group"></i>
             <span>二级目录管理</span>
         </div>
     </div>

     <div class="sidebar-content">
         <!--<div class="search-box">
             <i class="fas fa-search search-icon"></i>
             <input type="text" id="searchCollections" placeholder="搜索目录或请求...">
         </div> -->

         <div class="section-title">
             <h3>目录结构</h3>
             <div style="display: flex;">
                 <button class="add-btn" id="addFolder" >
                     <i class="fas fa-folder-plus"></i>
                 </button>
                 <button class="add-btn" id="addGlobalFeild" style="margin-left:10px;">
                     <i class="fas fa-file-alt"></i>
                 </button>
             </div>

         </div>

         <ul class="collection-list" id="collectionList">
             <!-- 目录和请求将通过JS动态生成 -->
         </ul>
     </div>
 </div>

 <!-- 主内容区域 -->
 <div class="main-content">
     <div class="header">
         <div class="app-title">
             <i class="fas fa-bolt"></i>
             <span>WebSocket 测试工具</span>
         </div>
         <div class="connection-status">
             <div class="status-indicator" id="statusIndicator"></div>
             <span id="statusText">未连接</span>
         </div>
         <div class="toolbar">
             <button class="toolbar-btn" id="connectBtn">
                 <i class="fas fa-plug"></i>
                 <span>连接</span>
             </button>
             <button class="toolbar-btn disconnect" id="disconnectBtn">
                 <i class="fas fa-power-off"></i>
                 <span>断开</span>
             </button>
         </div>
     </div>

     <div class="content-container">
         <!-- 连接配置面板 -->
         <div class="panel" style="height:25%">
             <div class="panel-header">
                 <span>连接配置</span>
                 <button class="btn-block" style="width: auto; padding: 5px 15px;" id="saveConfigBtn">
                     <i class="fas fa-save"></i> 保存配置
                 </button>
             </div>
             <div class="panel-body">
                 <div class="connection-form">
                     <div class="form-group">
                         <label for="connectionName">配置名称</label>
                         <input type="text" id="connectionName" class="form-control" placeholder="请求名称" value="">
                     </div>

                     <div class="form-group">
                         <label for="wsUrl">WebSocket URL</label>
                         <input type="text" id="wsUrl" class="form-control" placeholder="wss://example.com/socket" value="">
                     </div>

                     <div class="form-group">
                         <label for="protocol">协议</label>
                         <input type="text" id="protocol" class="form-control" placeholder="可选">
                     </div>

                     <div class="form-group">
                         <label for="timeout">超时时间(ms)</label>
                         <input type="number" id="timeout" class="form-control" value="5000">
                     </div>
                 </div>

                 <div class="form-group">
                     <label>请求头</label>
                     <div id="headersContainer">
                         <div class="form-row" style="display: flex; gap: 10px; margin-bottom: 10px;">
                             <input type="text" class="form-control" value="Authorization" style="flex: 1;">
                             <input type="text" class="form-control" value="Bearer token" style="flex: 1;">
                             <button class="action-btn" style="flex: 0 0 40px;">
                                 <i class="fas fa-times"></i>
                             </button>
                         </div>
                     </div>
                     <button class="btn-block" id="addHeaderBtn" style="margin-top: 10px;">
                         <i class="fas fa-plus"></i> 添加请求头
                     </button>
                 </div>

                 <!--<div class="form-group">
                     <label>初始消息</label>
                     <textarea id="initMessage" class="form-control" rows="3" placeholder='{"action":"subscribe","channel":"live"}'></textarea>
                 </div>-->
             </div>
         </div>

         <!-- 消息面板 -->
         <div class="panel" style="height:75%">
             <div class="panel-header">
                 <span>消息通信</span>
                 <div style="display: flex; gap: 10px;">
                     <button class="btn-block" style="width: auto; padding: 5px 15px;" id="historyMessagesBtn">
                         <i class="fas fa-history"></i> 历史
                     </button>
                     <button class="btn-block" style="width: auto; padding: 5px 15px;" id="clearMessages">
                         <i class="fas fa-trash"></i> 清空消息
                     </button>
                 </div>
             </div>
             <div class="panel-body">
                 <div class="messages-container">
                     <div class="message-history" id="messageHistory">
                         <div class="empty-state" id="emptyMessageState">
                             <div>📨</div>
                             <h3>没有消息</h3>
                             <p>连接后发送消息开始通信</p>
                         </div>
                     </div>

                     <div class="message-input-area">
                         <div class="input-options">
                             <button class="input-option-btn" id="formatJsonBtn" title="格式化JSON">
                                 <i class="fas fa-code"></i>
                             </button>
                         </div>
                         <textarea id="messageInput" class="message-input" placeholder="输入要发送的消息..."></textarea>
                         <div class="send-controls">
                             <select id="messageType" class="form-control">
                                 <option value="text">文本</option>
                                 <option value="json">JSON</option>
                                 <option value="binary">二进制</option>
                             </select>
                             <button class="send-btn" id="sendBtn">
                                 <i class="fas fa-paper-plane"></i> 发送
                             </button>
                         </div>
                     </div>
                 </div>
             </div>
         </div>
     </div>
 </div>

 <!-- 添加全局参数 -->
 <div class="modal" id="globalFeildModal">
     <div class="modal-content">
         <div class="modal-header">
             <i class="fas fa-folder-plus"></i> 添加全局参数
         </div>
         <div class="modal-body">
             <!--<div class="form-group">
                 <label for="folderName">参数名</label>
                 <input type="text" class="paramName" class="form-control" placeholder="输入参数名">
                 <label for="folderName">参数值</label>
                 <input type="text" class="paramValue" class="form-control" placeholder="输入值">
             </div>-->
         </div>
         <div class="modal-footer">
             <button class="toolbar-btn" id="cancelGlobalFeild">
                 <i class="fas fa-times"></i> 取消
             </button>
             <button class="toolbar-btn" style="background: var(--folder-color);" id="saveGlobalFeild">
                 <i class="fas fa-save"></i> 保存参数
             </button>
         </div>
     </div>
 </div>

 <!-- 添加目录模态框 -->
 <div class="modal" id="folderModal">
     <div class="modal-content">
         <div class="modal-header">
             <i class="fas fa-folder-plus"></i> 添加新目录
         </div>
         <div class="modal-body">
             <div class="form-group">
                 <label for="folderName">目录名称</label>
                 <input type="text" id="folderName" class="form-control" placeholder="输入目录名称">
             </div>
         </div>
         <div class="modal-footer">
             <button class="toolbar-btn" id="cancelFolder">
                 <i class="fas fa-times"></i> 取消
             </button>
             <button class="toolbar-btn" style="background: var(--folder-color);" id="saveFolder">
                 <i class="fas fa-save"></i> 保存目录
             </button>
         </div>
     </div>
 </div>

 <!-- 添加请求模态框 -->
 <div class="modal" id="requestModal">
     <div class="modal-content">
         <div class="modal-header">
             <i class="fas fa-plug"></i> 添加新请求
         </div>
         <div class="modal-body">
             <div class="form-group">
                 <label for="requestName">请求名称</label>
                 <input type="text" id="requestName" class="form-control" placeholder="输入请求名称">
             </div>
             <div class="form-group">
                 <label for="requestUrl">WebSocket URL</label>
                 <input type="text" id="requestUrl" class="form-control" placeholder="wss://example.com/socket">
             </div>
             <div class="form-group">
                 <label>父目录</label>
                 <select id="parentFolder" class="form-control">
                     <option value="">无 (顶级目录)</option>
                 </select>
             </div>
         </div>
         <div class="modal-footer">
             <button class="toolbar-btn" id="cancelRequest">
                 <i class="fas fa-times"></i> 取消
             </button>
             <button class="toolbar-btn" style="background: var(--request-color);" id="saveRequest">
                 <i class="fas fa-save"></i> 保存请求
             </button>
         </div>
     </div>
 </div>

 <!-- 历史消息模态框 -->
 <div class="modal" id="historyModal">
     <div class="modal-content" style="width: 600px; max-width: 90%;">
         <div class="modal-header">
             <i class="fas fa-history"></i> 历史消息
         </div>
         <div class="modal-body" style="max-height: 400px; overflow-y: auto;">
             <div id="historyList">
                 <!-- 历史消息将通过JS动态生成 -->
             </div>
         </div>
         <div class="modal-footer">
             <button class="toolbar-btn" id="clearHistoryBtn" style="background: var(--danger-color);">
                 <i class="fas fa-trash"></i> 清空历史
             </button>
             <button class="toolbar-btn" id="closeHistoryBtn">
                 <i class="fas fa-times"></i> 关闭
             </button>
         </div>
     </div>
 </div>
</main>

<style>

</style>
