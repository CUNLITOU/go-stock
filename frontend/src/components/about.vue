<script setup>
import { MdPreview } from 'md-editor-v3';
import 'md-editor-v3/lib/preview.css';
import {h, computed, nextTick, onBeforeUnmount, onMounted, ref} from 'vue';
import {GetConfig, GetVersionInfo,GetEffectiveSponsorVip,GetUserManual} from "../../wailsjs/go/main/App";
import {NTree, useNotification} from "naive-ui";
const updateLog = ref('');
const versionInfo = ref('');
const icon = ref('https://raw.githubusercontent.com/CUNLITOU/go-stock/master/build/appicon.png');
const alipay =ref('https://github.com/CUNLITOU/go-stock/raw/master/build/screenshot/alipay.jpg')
const wxpay =ref('https://github.com/CUNLITOU/go-stock/raw/master/build/screenshot/wxpay.jpg')
const wxgzh =ref('https://github.com/CUNLITOU/go-stock/raw/dev/build/screenshot/%E6%89%AB%E7%A0%81_%E6%90%9C%E7%B4%A2%E8%81%94%E5%90%88%E4%BC%A0%E6%92%AD%E6%A0%B7%E5%BC%8F-%E7%99%BD%E8%89%B2%E7%89%88.png')
const notify = useNotification()
// VIP 展示与功能门控共用后端 GetEffectiveSponsorVip 的权威判定结果：
// 只展示「当前是否生效」，避免关于页显示 VIP2 但 K线分析/AI助手提示权限不足。
const vipLevel=ref(0);
const vipActive=ref(false);
const vipReason=ref("");
const vipStartTime=ref("");
const vipEndTime=ref("");
const showManual = ref(false)
const manualContent = ref('')
const manualId = 'manual-preview'
const darkTheme = ref(false)
const theme = computed(() => darkTheme.value ? 'dark' : 'light')
const manualScrollRef = ref(null)
const catalogList = ref([])
const iframeLoading = ref(true)
// 支付/赞助页地址：与广场服务同站点，由后端下发的广场地址派生，不写死域名
const payUrl = ref('')

const buildCatalogTree = (headings) => {
  if (!headings.length) return []
  const roots = []
  const stack = []
  for (const h of headings) {
    const node = { key: h.text, label: h.text, level: h.level, children: [] }
    while (stack.length && stack[stack.length - 1].level >= h.level) {
      stack.pop()
    }
    if (stack.length) {
      stack[stack.length - 1].children.push(node)
    } else {
      roots.push(node)
    }
    stack.push(node)
  }
  const clean = (nodes) => {
    for (const n of nodes) {
      if (n.children.length === 0) delete n.children
      else clean(n.children)
    }
  }
  clean(roots)
  return roots
}

const catalogTree = computed(() => buildCatalogTree(catalogList.value))

const onTreeSelect = (keys) => {
  if (keys.length) scrollToHeading(keys[0])
}

const slugifyHeading = (text) => {
  return text
    .trim()
    .replace(/[^\w\u4e00-\u9fff]/g, '-')
    .replace(/-+/g, '-')
    .replace(/^-|-$/g, '')
}

const extractCatalog = () => {
  if (!manualScrollRef.value) return
  const headings = manualScrollRef.value.querySelectorAll('h1, h2, h3, h4, h5, h6')
  catalogList.value = Array.from(headings).map(h => ({
    text: h.textContent?.trim() || '',
    level: parseInt(h.tagName.slice(1))
  }))
}

const scrollToHeading = (headingText) => {
  if (!manualScrollRef.value) return
  const container = manualScrollRef.value
  const headings = container.querySelectorAll('h1, h2, h3, h4, h5, h6')
  for (const h of headings) {
    const text = h.textContent?.trim()
    if (text === headingText) {
      const containerRect = container.getBoundingClientRect()
      const headingRect = h.getBoundingClientRect()
      container.scrollTop += headingRect.top - containerRect.top - 10
      return
    }
  }
}

const openManual = () => {
  if (!manualContent.value) {
    GetUserManual().then(res => {
      manualContent.value = res
      showManual.value = true
      nextTick(() => { setTimeout(extractCatalog, 500) })
    })
  } else {
    showManual.value = true
    nextTick(() => { setTimeout(extractCatalog, 300) })
  }
}

onMounted(() => {
  document.title = '关于软件';
  GetConfig().then(res => {
    darkTheme.value = res.darkTheme
    if (res.promptPlazaApiBase) {
      payUrl.value = res.promptPlazaApiBase.replace(/\/api\/?$/, '') + '/pay'
    }
  })
  GetVersionInfo().then((res) => {
    updateLog.value = res.content;
    versionInfo.value = res.version;
    icon.value = res.icon;
    alipay.value=res.alipay;
    wxpay.value=res.wxpay;
    wxgzh.value=res.wxgzh;

  });

  // 以「当前是否生效」为准展示，与 K线分析 / AI助手 等功能的门控同源
  GetEffectiveSponsorVip().then((res) => {
    const lvl = Number(res?.vipLevel ?? 0);
    vipLevel.value = Number.isNaN(lvl) ? 0 : lvl;
    vipActive.value = !!res?.active;
    vipReason.value = res?.reason ?? "";
    vipStartTime.value = res?.startTime ?? "";
    vipEndTime.value = res?.endTime ?? "";
  }).catch(() => {
    vipLevel.value = 0;
    vipActive.value = false;
    vipReason.value = "";
  });



})
onBeforeUnmount(() => {
  notify.destroyAll()
})
</script>

<template>
      <div class="about-page" style="--wails-draggable:no-drag">
        <!-- 顶部信息区 -->
        <n-card size="large" :bordered="false" class="hero-card">
          <div class="hero">
            <div class="hero-icon">
              <n-image width="72" :src="icon" :preview-disabled="true" />
            </div>
            <div class="hero-title">
              <n-gradient-text type="info" :size="34" class="app-name">go-stock</n-gradient-text>
              <n-tag v-if="versionInfo" :bordered="false" type="success" size="small" round>
                v{{versionInfo}}
              </n-tag>
              <n-tag v-if="vipLevel > 0" :bordered="false" :type="vipActive ? 'warning' : 'error'" size="small" round>
                VIP{{vipLevel}}{{ vipActive ? '' : ' 未生效' }}
              </n-tag>
            </div>
            <n-gradient-text v-if="vipLevel > 0" :type="vipActive ? 'warning' : 'error'" class="vip-expire">
              <template v-if="vipActive">VIP 到期时间：{{ vipEndTime }}</template>
              <template v-else>{{ vipReason || 'VIP 未生效' }}<template v-if="!vipActive && vipStartTime">（生效时间：{{ vipStartTime }}）</template></template>
            </n-gradient-text>
            <n-flex justify="center" :size="12" class="hero-actions">
              <n-button size="small" @click="openManual" type="success" tertiary round>
                <template #icon>📖</template>
                查看用户手册
              </n-button>
              <n-button size="small" tag="a" href="https://github.com/CUNLITOU/go-stock" target="_blank" type="default" tertiary round>
                <template #icon>⭐</template>
                GitHub
              </n-button>
            </n-flex>
          </div>
        </n-card>

        <!-- 内嵌支付/赞助页 -->
        <n-card size="large" :bordered="false" class="embed-card" content-style="padding: 0;">
          <template #header>
            <div class="embed-header">
              <span class="embed-title">💖 赞助与会员</span>
              <n-button v-if="payUrl" size="tiny" quaternary type="primary" tag="a" :href="payUrl" target="_blank">
                在浏览器中打开 ↗
              </n-button>
            </div>
          </template>
          <div class="iframe-wrapper">
            <n-spin v-if="iframeLoading" class="iframe-spin" size="large" />
            <iframe
              v-if="payUrl"
              :src="payUrl"
              class="pay-iframe"
              :class="{ 'is-loaded': !iframeLoading }"
              allow="payment; clipboard-write"
              referrerpolicy="no-referrer-when-downgrade"
              @load="iframeLoading = false"
            ></iframe>
          </div>
        </n-card>

        <n-modal
          v-model:show="showManual"
          preset="card"
          title="用户手册"
          style="width: 90vw; max-height: 90vh"
          :bordered="false"
          :segmented="{ content: true, footer: true }"
        >
          <div style="display: flex; max-height: 75vh;">
            <div v-if="catalogList.length" class="manual-catalog" style="width: 240px; min-width: 240px; border-right: 1px solid var(--n-border-color); padding: 8px 4px; overflow-y: auto;">
              <div style="font-weight: bold; margin-bottom: 8px; padding: 0 8px;">目录</div>
              <n-tree
                :data="catalogTree"
                :block-line="true"
                :block-node="true"
                :selectable="true"
                :cancelable="false"
                default-expand-all
                key-field="key"
                label-field="label"
                children-field="children"
                @update:selected-keys="onTreeSelect"
              />
            </div>
            <div ref="manualScrollRef" style="flex: 1; overflow-y: auto; padding: 0 16px;">
              <MdPreview style="text-align: left;" :id="manualId" v-model="manualContent" :theme="theme" :preview-theme="'github'" :md-heading-id="slugifyHeading" @onHtmlChanged="extractCatalog" />
            </div>
          </div>
        </n-modal>
      </div>
</template>

<style scoped>
/* 页面整体 */
.about-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
  max-width: 1080px;
  margin: 0 auto;
}

/* 顶部信息区 */
.hero-card {
  border-radius: 12px;
}

.hero {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  padding: 8px 0 4px;
}

.hero-icon {
  width: 76px;
  height: 76px;
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(24, 160, 88, 0.25);
}

.hero-title {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
  justify-content: center;
}

.app-name {
  font-weight: 600;
  letter-spacing: 1px;
}

.vip-expire {
  font-size: 13px;
}

.hero-actions {
  margin-top: 4px;
}

/* 内嵌页 */
.embed-card {
  border-radius: 12px;
  overflow: hidden;
}

.embed-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.embed-title {
  font-size: 15px;
  font-weight: 600;
}

.iframe-wrapper {
  position: relative;
  border-top: 1px solid rgba(128, 128, 128, 0.12);
}

.iframe-spin {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 1;
}

.pay-iframe {
  display: block;
  width: 100%;
  height: calc(100vh - 320px);
  min-height: 560px;
  border: none;
  background: transparent;
  opacity: 0;
  transition: opacity 0.4s ease;
}

.pay-iframe.is-loaded {
  opacity: 1;
}

h1, h2 {
  margin: 0;
  padding: 6px 0;
}

p {
  margin: 2px 0;
}

ul {
  list-style-type: disc;
  padding-left: 20px;
}

a {
  color: #18a058;
  text-decoration: none;
}

a:hover {
  text-decoration: underline;
}

.manual-catalog > div:hover {
  color: #18a058;
}

.manual-catalog :deep(.n-tree-node-content) {
  text-align: left;
  justify-content: flex-start;
}

.manual-catalog :deep(.n-tree-node) {
  text-align: left;
}
</style>
