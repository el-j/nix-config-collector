<script setup>
import { ref, computed } from 'vue'
import { useToast } from 'primevue/usetoast'
import Button from 'primevue/button'
import ProgressBar from 'primevue/progressbar'
import Panel from 'primevue/panel'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import TabView from 'primevue/tabview'
import TabPanel from 'primevue/tabpanel'
import Toast from 'primevue/toast'
import Stepper from 'primevue/stepper'
import StepperPanel from 'primevue/stepperpanel'
import Textarea from 'primevue/textarea'
import InputText from 'primevue/inputtext'
import Divider from 'primevue/divider'

import { Scan, Preview, Generate } from './wailsjs/go/gui/App'

const toast = useToast()

// State
const loading = ref(false)
const loadingMsg = ref('')
const scanData = ref(null)
const nixConfig = ref(null)
const outputDir = ref('~/.config/nixpkgs')
const activeStep = ref(0)

// Computed
const packages = computed(() => scanData.value?.packages ?? [])
const services = computed(() => (scanData.value?.services ?? []).filter(s => s.enabled))
const fonts = computed(() => scanData.value?.fonts ?? [])
const dotfiles = computed(() => scanData.value?.dotfiles ?? [])

const packageTagSeverity = (type) => {
  switch (type) {
    case 'brew-formula': return 'success'
    case 'brew-cask': return 'info'
    case 'mas': return 'warn'
    default: return 'secondary'
  }
}

// Actions
async function doScan() {
  loading.value = true
  loadingMsg.value = 'Scanning macOS system…'
  try {
    const raw = await Scan()
    scanData.value = JSON.parse(raw)
    activeStep.value = 1
    toast.add({ severity: 'success', summary: 'Scan complete', detail: `Found ${packages.value.length} packages`, life: 3000 })
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Scan failed', detail: String(e), life: 5000 })
  } finally {
    loading.value = false
  }
}

async function doPreview() {
  loading.value = true
  loadingMsg.value = 'Generating Nix configuration…'
  try {
    nixConfig.value = await Preview()
    activeStep.value = 2
    toast.add({ severity: 'success', summary: 'Config generated', detail: 'Review your Nix configuration below', life: 3000 })
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Generation failed', detail: String(e), life: 5000 })
  } finally {
    loading.value = false
  }
}

async function doWrite() {
  loading.value = true
  loadingMsg.value = `Writing to ${outputDir.value}…`
  try {
    await Generate(outputDir.value)
    toast.add({ severity: 'success', summary: 'Written!', detail: `Config written to ${outputDir.value}`, life: 4000 })
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Write failed', detail: String(e), life: 5000 })
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <Toast position="top-right" />

  <div class="flex flex-col h-screen bg-zinc-950">
    <!-- Header -->
    <header class="flex items-center gap-3 px-6 py-4 border-b border-zinc-800 bg-zinc-900">
      <span class="text-2xl">❄️</span>
      <div>
        <h1 class="text-lg font-bold text-white leading-none">nix-config-collector</h1>
        <p class="text-xs text-zinc-400 mt-0.5">macOS → Nix configuration generator</p>
      </div>
      <div class="ml-auto flex items-center gap-2">
        <Tag v-if="scanData" severity="success" icon="pi pi-check-circle" :value="`${packages.length} packages`" />
        <Tag v-if="nixConfig" severity="info" icon="pi pi-file" value="Config ready" />
      </div>
    </header>

    <!-- Loading overlay -->
    <div v-if="loading" class="fixed inset-0 z-50 bg-black/70 flex flex-col items-center justify-center gap-4">
      <div class="text-4xl animate-spin">❄️</div>
      <p class="text-white font-medium">{{ loadingMsg }}</p>
      <ProgressBar mode="indeterminate" class="w-64 h-1" />
    </div>

    <!-- Main content -->
    <main class="flex-1 overflow-auto p-6">
      <Stepper :value="activeStep" class="mb-6">
        <StepperPanel :value="0">
          <template #header="{ value, activateCallback }">
            <button @click="activateCallback" class="flex items-center gap-2 px-3 py-1.5">
              <span class="w-7 h-7 rounded-full flex items-center justify-center text-sm font-bold"
                :class="value <= activeStep ? 'bg-blue-500 text-white' : 'bg-zinc-700 text-zinc-400'">1</span>
              <span :class="value <= activeStep ? 'text-white' : 'text-zinc-400'" class="font-medium">Scan</span>
            </button>
          </template>
          <template #content>
            <div class="py-4 space-y-4">
              <p class="text-zinc-400">Scan your macOS system for installed packages, services, dotfiles, and preferences.</p>
              <Button @click="doScan" icon="pi pi-search" label="Scan System" :loading="loading" severity="primary" size="large" />
            </div>
          </template>
        </StepperPanel>

        <StepperPanel :value="1">
          <template #header="{ value, activateCallback }">
            <button @click="activateCallback" :disabled="!scanData" class="flex items-center gap-2 px-3 py-1.5">
              <span class="w-7 h-7 rounded-full flex items-center justify-center text-sm font-bold"
                :class="value <= activeStep ? 'bg-blue-500 text-white' : 'bg-zinc-700 text-zinc-400'">2</span>
              <span :class="value <= activeStep ? 'text-white' : 'text-zinc-400'" class="font-medium">Preview</span>
            </button>
          </template>
          <template #content>
            <div class="py-4 space-y-4">
              <!-- Scan summary -->
              <div v-if="scanData" class="grid grid-cols-2 md:grid-cols-4 gap-3 mb-4">
                <div v-for="stat in [
                  { label: 'Packages', value: packages.length, icon: 'pi-box' },
                  { label: 'Services', value: services.length, icon: 'pi-server' },
                  { label: 'Dotfiles', value: dotfiles.length, icon: 'pi-file' },
                  { label: 'Fonts', value: fonts.length, icon: 'pi-palette' },
                ]" :key="stat.label"
                  class="bg-zinc-800 rounded-lg p-3 flex items-center gap-3">
                  <i :class="`pi ${stat.icon} text-blue-400 text-xl`"></i>
                  <div>
                    <div class="text-2xl font-bold text-white">{{ stat.value }}</div>
                    <div class="text-xs text-zinc-400">{{ stat.label }}</div>
                  </div>
                </div>
              </div>

              <!-- Package table -->
              <Panel v-if="packages.length" header="Packages" :toggleable="true" class="mb-3">
                <DataTable :value="packages" :rows="10" paginator scrollable scrollHeight="300px"
                  size="small" stripedRows class="text-sm">
                  <Column field="name" header="Name" sortable />
                  <Column field="type" header="Type" sortable>
                    <template #body="{ data }">
                      <Tag :severity="packageTagSeverity(data.type)" :value="data.type" />
                    </template>
                  </Column>
                  <Column field="version" header="Version" />
                  <Column field="nix_name" header="Nix attr" />
                </DataTable>
              </Panel>

              <Button @click="doPreview" icon="pi pi-eye" label="Generate Preview" :loading="loading"
                :disabled="!scanData" severity="success" size="large" />
            </div>
          </template>
        </StepperPanel>

        <StepperPanel :value="2">
          <template #header="{ value, activateCallback }">
            <button @click="activateCallback" :disabled="!nixConfig" class="flex items-center gap-2 px-3 py-1.5">
              <span class="w-7 h-7 rounded-full flex items-center justify-center text-sm font-bold"
                :class="value <= activeStep ? 'bg-blue-500 text-white' : 'bg-zinc-700 text-zinc-400'">3</span>
              <span :class="value <= activeStep ? 'text-white' : 'text-zinc-400'" class="font-medium">Write</span>
            </button>
          </template>
          <template #content>
            <div class="py-4 space-y-4">
              <!-- Config preview tabs -->
              <TabView v-if="nixConfig">
                <TabPanel v-for="tab in [
                  { key: 'darwin', label: 'darwin-configuration.nix', content: nixConfig.DarwinConfig },
                  { key: 'home', label: 'home.nix', content: nixConfig.HomeConfig },
                  { key: 'flake', label: 'flake.nix', content: nixConfig.FlakeConfig },
                ]" :key="tab.key" :header="tab.label">
                  <pre class="code-block bg-zinc-900 rounded-lg p-4 overflow-auto max-h-64 text-green-300 text-xs">{{ tab.content }}</pre>
                </TabPanel>
              </TabView>

              <Divider />

              <!-- Output directory + write button -->
              <div class="flex items-center gap-3">
                <label class="text-zinc-300 font-medium whitespace-nowrap">Output directory:</label>
                <InputText v-model="outputDir" class="flex-1 font-mono" />
              </div>
              <Button @click="doWrite" icon="pi pi-save" label="Write Config Files" :loading="loading"
                :disabled="!nixConfig" severity="warn" size="large" />

              <p class="text-zinc-500 text-sm">
                After writing, apply with:
                <code class="text-green-400 font-mono text-xs bg-zinc-900 px-2 py-1 rounded ml-1">
                  darwin-rebuild switch --flake {{ outputDir }}#$(hostname)
                </code>
              </p>
            </div>
          </template>
        </StepperPanel>
      </Stepper>
    </main>

    <!-- Footer -->
    <footer class="px-6 py-3 border-t border-zinc-800 bg-zinc-900 flex items-center justify-between">
      <span class="text-xs text-zinc-500">nix-config-collector</span>
      <span v-if="scanData" class="text-xs text-zinc-500">
        Scanned: {{ new Date(scanData.scanned_at).toLocaleString() }} · {{ scanData.hostname }}
      </span>
    </footer>
  </div>
</template>
