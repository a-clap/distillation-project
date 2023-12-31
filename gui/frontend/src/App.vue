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

<template>
  <html class="dark">
  <div class="app">
    <Sidebar/>
    <el-dialog v-model="err.show" :title="err.title" width="70%" :modal=false :center=true :close-on-click-modal=false
               :show-close=false align-center>
      <span class="dialog-message">{{ err.msg }}</span>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="success" size="large" @click="err.close">
            {{ $t('errors.submit') }}
          </el-button>
          <el-button type="danger" size="large" @click="err.skip">
            {{ $t('errors.skip') }}
          </el-button>
        </span>
      </template>
    </el-dialog>
    <el-dialog v-model="updater.commit" :title="$t('system.update_success')" width="70%" :modal=false :center=true
               :close-on-click-modal=false
               :show-close=true align-center>
      <span class="dialog-message"> {{ updater.message }} </span>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="success" size="large" @click="updater.submit">
            {{ $t('errors.submit') }}
          </el-button>
        </span>
      </template>
    </el-dialog>
    <router-view/>
  </div>

  </html>
</template>

<script setup lang="ts">

import Sidebar from "./components/Sidebar.vue";
import {useGpioStore} from "./stores/gpios";
import {useDSStore} from "./stores/ds";
import {useHeatersStore} from "./stores/heaters";
import {usePTStore} from "./stores/pt";
import {useWIFIStore} from "./stores/wifi";
import {usePhasesStore} from "./stores/phases";
import {useErrorStore} from "./stores/errors";
import {useLogStore} from "./stores/log";
import {useProcessStore} from "./stores/process";
import {useNameStore} from "./stores/names";
import {useUpdaterStore} from "./stores/updater";

const err = useErrorStore()
const updater = useUpdaterStore()

interface StoreInitializer {
  init: Function;
}

let initFuncs: StoreInitializer[] = [
  useErrorStore(),
  useProcessStore(),
  useGpioStore(),
  useDSStore(),
  useHeatersStore(),
  usePTStore(),
  useWIFIStore(),
  usePhasesStore(),
  useLogStore(),
  useNameStore(),
  useUpdaterStore(),
]

initFuncs.forEach((store) => {
  setTimeout(() => {
    store.init()
  }, 10);
})

</script>
<style lang="scss">
:root {
  --sidebar-dark: #1e293b;
  --sidebar-dark-alt: #334155;
  --sidebar-width: 200px;
  --window-width: 1024px;
  --window-height: 768px;
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
  font-family: 'Fira sans', sans-serif;
  transition: none !important;
  transform: none !important;
}

.app {
  display: flex;
  background-color: rgb(51, 65, 85);

  main {
    flex: 1 1 0;
    padding-left: 0.5rem;
    padding-top: 0.5rem;
  }

  input {
    width: 100px;
    height: 34px;
    padding: 6px 12px;
    line-height: 1rem;
    text-align: center;
    color: #555;
    cursor: default;
    caret-color: transparent;
    background-color: #fff;
    background-image: none;
    border: 1px solid #ccc;
    border-radius: 4px;
    box-shadow: inset 0 1px 1px rgba(0, 0, 0, .075);
  }

  input:focus {
    outline: none !important;
    border: 1px solid var(--el-color-primary);
  }

  .dialog-message {
    display: flex;
    justify-content: space-around;
    font-size: 1.5rem;

  }

  .dialog-footer {
    display: flex;
    justify-content: space-around;
  }
}
</style>
