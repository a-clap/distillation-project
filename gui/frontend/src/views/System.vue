<template>
  <main>
    <h1>{{ $t('system.title') }}</h1>
    <el-tabs v-model="activated" type="card" tabPosition="left" class="demo-tabs">
      <el-tab-pane :label="$t('system.names')" name="names" class="main-tab">
        <div class="name-container">{{ $t('system.ds') }}</div>
        <template v-for="(ds, index) in dsStore.ds" :key="index">
          <section class="name-container">
            <div class="name-elem">{{ ds.id }}</div>
            <input class="name-input" v-model="ds.name.view" @click="() => ds.name.showKeyboard()">
            <Keyboard v-bind="ds.name" :write="(e: string) => ds.writeName(e)" :cancel="() => ds.name.cancel()"/>
          </section>
        </template>
        <div class="name-container">{{ $t('system.pt') }}</div>
        <template v-for="(pt, index) in ptStore.pt" :key="index">
          <section class="name-container">
            <div class="name-elem">{{ pt.id }}</div>
            <input class="name-input" v-model="pt.name.view" @click="() => pt.name.showKeyboard()">
            <Keyboard v-bind="pt.name" :write="(e: string) => pt.writeName(e)" :cancel="() => pt.name.cancel()"/>
          </section>
        </template>
      </el-tab-pane>
      <el-tab-pane :label="$t('system.time')" name="time">
        <div class="ntp-container">
          <div>{{ $t('system.time_now') }}: {{ timeNow }}</div>
          <el-switch v-model="ntp" :active-text="$t('system.ntp_time')" size="large"/>
        </div>

        <div class="time-picker" v-if=!ntpEnabled>
          <el-time-picker
              v-model="currentTime"
              :placeholder="$t('system.pick_time')"
              :clearable=false
              size="large"
          />
          <el-date-picker
              v-model="currentDate"
              type="date"
              :placeholder="$t('system.pick_date')"
              :clearable=false
              size="large"
          />
          <el-button type="primary" size="large" @click="setTime">
            {{ $t('system.set_time') }}
          </el-button>
        </div>
      </el-tab-pane>
      <el-tab-pane :label="$t('system.net')" name="net">
        <section class="net-interface">
          <h3>{{ $t(('system.net_interface')) }}</h3>
          <h3>{{ $t(('system.net_ip')) }}</h3>
        </section>
        <section v-for="(netInterface, index) in netInterfaces" :key="index">
          <section class="net-interface">
            <div> {{ netInterface.name }}</div>
            <div> {{ netInterface.ip_addr }}</div>
          </section>
        </section>
      </el-tab-pane>
      <el-tab-pane :label="$t('system.update')" name="update">
        <section class="update-interface">
          <el-button type="primary" :icon="Refresh" size="large"
                     @click=checkUpdate>
            Sprawdz aktualizacje
          </el-button>
          <h3 v-if="updaterStore.updating">Pobieranie</h3>
          <el-progress :stroke-width=20 :percentage=progress :color="colors" :format="(p: number) => {return p +`%`}"/>
          <h3>Instalowanie</h3>
          <el-progress :stroke-width=20 :percentage=progress :color="colors" :format="(p: number) => {return p +`%`}"/>
        </section>
      </el-tab-pane>
    </el-tabs>
  </main>
</template>

<script setup lang="ts">

import {useDSStore} from "../stores/ds";
import {usePTStore} from "../stores/pt";
import {computed, onMounted, ref} from "vue";
import Keyboard from "../components/Keyboard.vue";
import {ListInterfaces, Now, NTPGet, NTPSet, TimeSet} from "../../wailsjs/go/backend/Backend";
import {backend} from "../../wailsjs/go/models";
import {Loader} from "../types/Loader";
import {AppErrorCodes} from "../stores/error_codes";
import {i18n} from "../i18n";
import {FormatDate} from "../stores/log";
import {Refresh} from "@element-plus/icons-vue";
import {useUpdaterStore} from "../stores/updater";
import NetInterface = backend.NetInterface;

const currentDate = ref(new Date())
const currentTime = ref(new Date())
const activated = ref('update')
const dsStore = useDSStore()
const ptStore = usePTStore()
const updaterStore = useUpdaterStore()

const ntpEnabled = ref(false)
const ntp = computed({
  get: () => ntpEnabled.value,
  set: (v: boolean) => {

    let msg = i18n.global.t('system.ntp_loading')
    Loader.show(AppErrorCodes.NTPFailed, 5000, msg)
    if (!v) {
      currentTime.value = new Date()
    }

    NTPSet(v).then((err: any) => {
      if (!err) {
        ntpEnabled.value = v
      }
      Loader.close()
    })
  }
})

const netInterfaces = ref<NetInterface[]>([])
const timeNow = ref('')
const progress = ref(0)

onMounted(() => {

  NTPGet().then((value: boolean) => {
    ntp.value = value
  })

  setInterval(() => {
    Now().then((ts: number) => {
      timeNow.value = FormatDate(new Date(ts))
    })
  }, 1000)

  ListInterfaces().then((interfaces: NetInterface[]) => {
    netInterfaces.value = interfaces
  })

  setInterval(() => {
    progress.value = progress.value >= 100 ? 0 : progress.value + 1;
  }, 100)

})

function checkUpdate() {
  Loader.show(AppErrorCodes.NTPFailed, 5000, "sprawdam update")

  updaterStore.checkUpdate()
  Loader.close()
}

function setTime() {
  let fullDate = currentTime.value
  fullDate.setFullYear(currentDate.value.getFullYear())
  fullDate.setMonth(currentDate.value.getMonth())
  fullDate.setDate(currentDate.value.getDate())

  TimeSet(fullDate.getTime()).then(() => {

  })
}

const percentage = ref(0)

const colors = [
  {color: '#FF0000', percentage: 25},
  {color: '#FF7F00', percentage: 50},
  {color: '#FFFF00', percentage: 75},
  {color: '#00FF00', percentage: 100},
]


</script>

<style lang="scss" scoped>
h1 {
  margin-bottom: 1rem;
}

.demo-tabs > .el-tabs__content {
  padding: 32px;
}

.name-input {
  width: 150px
}

.name-container {
  display: flex;
  justify-content: center;
  align-items: center;
  margin-bottom: 2rem;
  font-size: 2rem;

  .name-input {
    margin-left: 3rem;
  }

  .keyboard-window {
    font-size: initial;
  }
}

.ntp-container {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

.ntp-container > * {
  margin-bottom: 2rem;
}

.net-interface {
  display: flex;
  justify-content: space-around;
  align-items: center;

  * {
    margin-bottom: 1rem;
  }
}

.time-picker {
  display: flex;
  flex-direction: row;
  justify-content: space-around;
  align-items: center;

}

.update-interface {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;

  * {
    margin-bottom: 2rem;
  }
}

.el-progress--line {
  width: 400px;
}

.el-progress {

}

</style>
