<!-- MIT License -->
<!-- -->
<!-- Copyright (c) 2023 a-clap -->
<!-- -->
<!-- Permission is hereby granted, free of charge, to any person obtaining a copy -->
<!-- of this software and associated documentation files (the "Software"), to deal -->
<!-- in the Software without restriction, including without limitation the rights -->
<!-- to use, copy, modify, merge, publish, distribute, sublicense, and/or sell -->
<!-- copies of the Software, and to permit persons to whom the Software is -->
<!-- furnished to do so, subject to the following conditions: -->
<!-- -->
<!-- The above copyright notice and this permission notice shall be included in all -->
<!-- copies or substantial portions of the Software. -->
<!-- -->
<!-- THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR -->
<!-- IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, -->
<!-- FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE -->
<!-- AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER -->
<!-- LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, -->
<!-- OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE -->
<!-- SOFTWARE. -->

<template xmlns="http://www.w3.org/1999/html">
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
          <el-button :disabled=updaterStore.updating type="primary" :icon="Refresh" size="large" @click=checkUpdate>
            {{ $t(('system.check_update')) }}
          </el-button>
          <h3 v-if="updaterStore.new_update"> {{ $t(('system.release_name')) + release }} </h3>
          <section v-if="updaterStore.new_update">
            <el-button :disabled=updaterStore.updating type="success" :icon="Download" size="large"
                       @click=updaterStore.startUpdate()>
              {{ $t(('system.start_update')) }}
            </el-button>
            <el-button :disabled=!updaterStore.updating type="danger" :icon="CircleClose" size="large"
                       @click=updaterStore.stopUpdate()>
              {{ $t(('system.stop_update')) }}
            </el-button>
          </section>
          <section v-if="updaterStore.updating" class="bars">
            <h3>{{ $t(('system.downloading')) }}</h3>
            <el-progress :stroke-width=20 :percentage=updaterStore.downloading :color="colors"
                         :format="(p: number) => {return p +`%`}"/>

            <h3>{{ $t(('system.installing')) }}</h3>
            <el-progress :stroke-width=20 :percentage=updaterStore.installing :color="colors"
                         :format="(p: number) => {return p +`%`}"/>

            <h3>{{ $t(('system.rebooting')) }}</h3>
            <el-progress :stroke-width=20 :percentage=updaterStore.rebooting :color="colors"
                         :format="(p: number) => {return p +`%`}"/>
          </section>
        </section>
      </el-tab-pane>
    </el-tabs>
  </main>
</template>

<script setup lang="ts">

import {useDSStore} from "../stores/ds";
import {usePTStore} from "../stores/pt";
import {computed, onMounted, onUnmounted, ref} from "vue";
import Keyboard from "../components/Keyboard.vue";
import {ListInterfaces, Now, NTPGet, NTPSet, TimeSet} from "../../wailsjs/go/backend/Backend";
import {backend} from "../../wailsjs/go/models";
import {Loader} from "../types/Loader";
import {AppErrorCodes} from "../stores/error_codes";
import {i18n} from "../i18n";
import {FormatDate} from "../stores/log";
import {CircleClose, Download, Refresh} from "@element-plus/icons-vue";
import {useUpdaterStore} from "../stores/updater";
import NetInterface = backend.NetInterface;

const currentDate = ref(new Date())
const currentTime = ref(new Date())
const activated = ref('names')
const dsStore = useDSStore()
const ptStore = usePTStore()
const updaterStore = useUpdaterStore()

const release = ref('')

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
const timer = ref(0)

onMounted(() => {
  timer.value = setInterval(() => {
    Now().then((ts: number) => {
      timeNow.value = FormatDate(new Date(ts))
    })
  }, 1000)

  NTPGet().then((value: boolean) => {
    ntp.value = value
  })


  ListInterfaces().then((interfaces: NetInterface[]) => {
    netInterfaces.value = interfaces
  })
})

onUnmounted(() => {
  clearInterval(timer.value)
})

function checkUpdate() {
  let msg = i18n.global.t('system.pulling_updates')
  Loader.show(AppErrorCodes.CheckUpdates, 5000, msg)

  updaterStore.checkUpdate().then(() => {
    if (updaterStore.releases.length > 0) {
      release.value = updaterStore.releases[0]
    }

    Loader.close()
  }).catch(function (error) {
    if (error === parseInt(error, 10)) {
      Loader.failNow(error)
    }
    console.error(error)
  })
}


function setTime() {

  let fullDate = currentTime.value
  // Construct proper time
  fullDate.setFullYear(currentDate.value.getFullYear())
  fullDate.setMonth(currentDate.value.getMonth())
  fullDate.setDate(currentDate.value.getDate())

  TimeSet(fullDate.getTime()).catch(function (error) {
    console.log("getTime" + error)
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

.update-interface, .bars {
  display: flex;
  flex-direction: column;
  align-items: center;

  * {
    margin-bottom: 1.5rem;
  }
}


.el-progress--line {
  width: 400px;
}

.el-progress {

}

</style>
