<template>
  <main class="phases-page">
    <h1>{{ $t('phases.title') }}</h1>
    <el-tabs v-model="activated" type="card" tabPosition="left" class="demo-tabs">
      <el-tab-pane :label="$t('phases.main')" name="main" class="main-tab">
        <el-row :gutter="20" align="middle">
          <el-col :span="7" :offset="3">
            <el-button type="primary" size="large" @click="() => LoadParameters()">{{ $t('phases.load') }}</el-button>
          </el-col>
          <el-col :span="7" :offset="3">
            <el-button type="primary" size="large" @click="() => SaveParameters()">{{ $t('phases.save') }}</el-button>
          </el-col>
        </el-row>
        <el-row :gutter="20" align="middle">
          <el-col :span="5" :offset="5">
            <label>{{ $t('phases.count') }}</label>
          </el-col>
          <el-col :span="7" :offset="3">
            <input v-model="phaseStore.phases.phaseCount.view"
                   @click="() => phaseStore.phases.phaseCount.showKeyboard()">
            <Keyboard v-bind="phaseStore.phases.phaseCount"
                      :write="(e: number) => phaseStore.phases.phaseCount.write(e)"
                      :cancel="() => phaseStore.phases.phaseCount.cancel()"/>
          </el-col>
        </el-row>
        <el-row :gutter="20" align="middle">
          <el-col :span="15" :offset="6" :style="'font-weight: 800'">
            <label v-if="phaseStore.phases.gpios.length == 0">{{ $t('phases.no_gpio') }}</label>
            <label v-else>{{ $t('phases.gpio_global') }}</label>
          </el-col>
        </el-row>
        <template v-for="( gpio ) in  phaseStore.phases.gpios ">
          <el-row :gutter="20">
            <el-col :span="5">
              <el-switch v-model="gpio.enable" :active-text="gpio.id" size="large"/>
            </el-col>
            <el-col :span="4" v-if="gpio.enable">
              <el-checkbox v-model="gpio.inverted" :label="$t('phases.gpio_inverted')" size="large" border/>
            </el-col>
          </el-row>
          <el-row :gutter="20" v-if="gpio.enable">
            <el-col :span="4">
              <label>{{ $t('phases.next_sensor') }}</label>
            </el-col>
            <el-col :span="6">
              <el-select v-model="gpio.sensor_id" size="large" class="m-2">
                <el-option v-for="sensor in phaseStore.phases.sensors" :label="sensor" :value="sensor"/>
              </el-select>
            </el-col>
            <el-col :span="5" :offset="2">
              <label>{{ $t('phases.gpio_hysteresis') }}</label>
            </el-col>
            <el-col :span="3" :offset="1">
              <input v-model="gpio.hysteresis.view" @click="() => gpio.hysteresis.showKeyboard()">
              <Keyboard v-bind="gpio.hysteresis" :write="(e: number) => gpio.hysteresis.write(e)"
                        :cancel="() => gpio.hysteresis.cancel()"/>
            </el-col>
          </el-row>
          <el-row :gutter="20" v-if="gpio.enable">
            <el-col :span="7">
              <label>{{ $t('phases.gpio_temp_min') }}</label>
            </el-col>
            <el-col :span="3">
              <input v-model="gpio.t_low.view" @click="() => gpio.t_low.showKeyboard()">
              <Keyboard v-bind="gpio.t_low" :write="(e: number) => gpio.t_low.write(e)"
                        :cancel="() => gpio.t_low.cancel()"/>
            </el-col>
            <el-col :span="7" :offset=1>
              <label>{{ $t('phases.gpio_temp_max') }}</label>
            </el-col>
            <el-col :span="3">
              <input v-model="gpio.t_high.view" @click="() => gpio.t_high.showKeyboard()">
              <Keyboard v-bind="gpio.t_high" :write="(e: number) => gpio.t_high.write(e)"
                        :cancel="() => gpio.t_high.cancel()"/>
            </el-col>
          </el-row>
        </template>
      </el-tab-pane>
      <el-tab-pane v-for="(phase, index) in  phaseStore.phases.phases " :key="index"
                   :label="$t('phases.phase') + ' ' + (index + 1)">
        <el-row :gutter="20" align="middle">
          <el-col :span="5">
            {{ $t('phases.next_type') }}
          </el-col>
          <el-col :span="7">
            <el-switch v-model="phase.next_type" :inactive-text="$t('phases.next_type_time')"
                       :active-text="$t('phases.next_type_temp')" size="large"/>
          </el-col>
          <el-col :span="7" :offset="1">
            <label v-if="!phase.next_type">{{ $t('phases.timeleft') }}</label>
            <label v-else>{{ $t('phases.timeleft_temp') }}</label>
          </el-col>
          <el-col :span="3">
            <input v-model="phase.next_timeleft.view" @click="() => phase.next_timeleft.showKeyboard()">
            <Keyboard v-bind="phase.next_timeleft" :write="(e: number) => phase.next_timeleft.write(e)"
                      :cancel="() => phase.next_timeleft.cancel()"/>
          </el-col>
        </el-row>
        <el-row :gutter="20" align="middle" v-if="phase.next_type">
          <el-col :span="4">
            <label>{{ $t('phases.next_sensor') }}</label>
          </el-col>
          <el-col :span="5">
            <el-select v-model="phase.next_sensor" size="large" class="m-2">
              <el-option v-for="sensor in phaseStore.phases.sensors" :label="sensor" :value="sensor"/>
            </el-select>
          </el-col>
          <el-col :span="6" :offset="4">
            <label v-if="phase.next_type">{{ $t('phases.sensor_threshold') }}</label>
          </el-col>
          <el-col :span="3" :offset="1">
            <input v-model="phase.next_sensor_threshold.view" @click="() => phase.next_sensor_threshold.showKeyboard()">
            <Keyboard v-bind="phase.next_sensor_threshold" :write="(e: number) => phase.next_sensor_threshold.write(e)"
                      :cancel="() => phase.next_sensor_threshold.cancel()"/>
          </el-col>
        </el-row>
        <el-row :gutter="20" align="middle">
          <el-col :span="15" :offset="8" :style="'font-weight: 800'">
            <label v-if="phase.heaters.length == 0">{{ $t('phases.no_heaters') }}</label>
            <label v-else>{{ $t('phases.heaters') }}</label>
          </el-col>
        </el-row>
        <el-row :gutter="20" v-for="(heater) in  phase.heaters ">
          <el-col :span="5" :style="'font-weight: 800'">
            <label>{{ heater.id }}</label>
          </el-col>
          <el-col :span="4">
            <input v-model="heater.power.view" @click="() => heater.power.showKeyboard()">
            <Keyboard v-bind="heater.power" :write="(e: number) => heater.power.write(e)"
                      :cancel="() => heater.power.cancel()"/>
          </el-col>
          <el-col :span="5">
            <label>{{ $t('phases.heater_power') }}</label>
          </el-col>
        </el-row>
        <el-row :gutter="20" align="middle">
          <el-col :span="15" :offset="8" :style="'font-weight: 800'">
            <label v-if="phase.gpios.length == 0">{{ $t('phases.no_gpio') }}</label>
            <label v-else>{{ $t('phases.gpio') }}</label>
          </el-col>
        </el-row>
        <template v-for="( gpio ) in  phase.gpios ">
          <el-row :gutter="20">
            <el-col :span="5">
              <el-switch v-model="gpio.enable" :active-text="gpio.id" size="large"/>
            </el-col>
            <el-col :span="4" v-if="gpio.enable">
              <el-checkbox v-model="gpio.inverted" :label="$t('phases.gpio_inverted')" size="large" border/>
            </el-col>
          </el-row>
          <el-row :gutter="20" v-if="gpio.enable">
            <el-col :span="4">
              <label>{{ $t('phases.next_sensor') }}</label>
            </el-col>
            <el-col :span="6">
              <el-select v-model="gpio.sensor_id" size="large" class="m-2">
                <el-option v-for="sensor in phaseStore.phases.sensors" :label="sensor" :value="sensor"/>
              </el-select>
            </el-col>
            <el-col :span="5" :offset="2">
              <label>{{ $t('phases.gpio_hysteresis') }}</label>
            </el-col>
            <el-col :span="3" :offset="1">
              <input v-model="gpio.hysteresis.view" @click="() => gpio.hysteresis.showKeyboard()">
              <Keyboard v-bind="gpio.hysteresis" :write="(e: number) => gpio.hysteresis.write(e)"
                        :cancel="() => gpio.hysteresis.cancel()"/>
            </el-col>
          </el-row>
          <el-row :gutter="20" v-if="gpio.enable">
            <el-col :span="7">
              <label>{{ $t('phases.gpio_temp_min') }}</label>
            </el-col>
            <el-col :span="3">
              <input v-model="gpio.t_low.view" @click="() => gpio.t_low.showKeyboard()">
              <Keyboard v-bind="gpio.t_low" :write="(e: number) => gpio.t_low.write(e)"
                        :cancel="() => gpio.t_low.cancel()"/>
            </el-col>
            <el-col :span="7" :offset=1>
              <label>{{ $t('phases.gpio_temp_max') }}</label>
            </el-col>
            <el-col :span="3">
              <input v-model="gpio.t_high.view" @click="() => gpio.t_high.showKeyboard()">
              <Keyboard v-bind="gpio.t_high" :write="(e: number) => gpio.t_high.write(e)"
                        :cancel="() => gpio.t_high.cancel()"/>
            </el-col>
          </el-row>
        </template>
      </el-tab-pane>
    </el-tabs>
  </main>
</template>

<script setup lang="ts">

import Keyboard from "../components/Keyboard.vue"
import {onMounted, ref} from "vue";
import {usePhasesStore} from "../stores/phases";
import {LoadParameters, SaveParameters} from "../../wailsjs/go/backend/Backend";

const activated = ref('main')
const phaseStore = usePhasesStore()

onMounted(() => {
  phaseStore.reload()
})

</script>

<style lang="scss">
.demo-tabs > .el-tabs__content {
  padding: 32px;
}

.el-row {
  margin-bottom: 1rem;
  align-items: center;
}
</style>
