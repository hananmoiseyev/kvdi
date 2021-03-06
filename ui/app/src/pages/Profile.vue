<!--
Copyright 2020,2021 Avi Zimmerman

This file is part of kvdi.

kvdi is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

kvdi is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with kvdi.  If not, see <https://www.gnu.org/licenses/>.
-->

<template>
  <q-page padding>
    <div class="wrapper">

      <!-- Header -->
      <div class="q-pa-md row items-start q-gutter-md">
        <div class="text-h4">User Settings</div>
      </div>

      <div class="q-pa-md row items-start q-gutter-md">

        <!-- Password options -->
        <q-card class="bg-grey-1" style="width:500px">
          <q-card-section>
            <div class="row items-center no-wrap">
              <div class="text-h6"><q-icon name="security" />&nbsp;Change Password</div>
            </div>
          </q-card-section>
          <q-card-section>
            <PasswordInput ref="password" :startDisabled="true" />
            <q-btn :disabled="passwordSubmitDisabled" color="primary" flat label="Cancel" @click="resetPasswordInput" />
            <q-btn :disabled="passwordSubmitDisabled" color="primary" flat label="Update" @click="doUpdatePassword" />
          </q-card-section>
        </q-card>

      </div>

      <div class="q-pa-md row items-start q-gutter-md">
        <!-- MFA Config -->
        <q-card class="bg-grey-1" style="width:500px">
          <q-card-section>
            <div class="row items-center no-wrap">
              <div class="text-h6"><q-icon name="security" />&nbsp;MFA Options</div>
            </div>
          </q-card-section>
          <q-card-section>
            <MFAConfig ref="mfaconfig" :username="username" :newUser="false" />
          </q-card-section>
        </q-card>
      </div>

    </div>
  </q-page>
</template>

<script>
import PasswordInput from 'components/inputs/PasswordInput.vue'
import MFAConfig from 'components/inputs/MFAConfig.vue'

export default {
  name: 'Profile',
  components: { PasswordInput, MFAConfig },
  mounted () { this.$refs.password.password = '*****************************' },
  created () { this.$root.$on('edit-password', this.setEditPassword) },
  beforeDestroy () { this.$root.$off('edit-password', this.setEditPassword) },
  data () {
    return {
      passwordSubmitDisabled: true
    }
  },
  computed: {
    username () {
      return this.$userStore.getters.user.name
    }
  },
  methods: {
    resetPasswordInput () {
      this.$refs.password.passwordIsDisabled = true
      this.passwordSubmitDisabled = true
      this.$refs.password.password = '*****************************'
    },
    setEditPassword () {
      this.passwordSubmitDisabled = false
    },
    async doUpdatePassword () {
      if (this.$refs.password.passwordIsDisabled) { return }
      const payload = {
        password: this.$refs.password.password
      }
      const user = this.username
      try {
        await this.$axios.put(`/api/users/${user}`, payload)
        this.$q.notify({
          color: 'green-4',
          textColor: 'white',
          icon: 'cloud_done',
          message: `Password updated successfully for '${user}'`
        })
        this.resetPasswordInput()
      } catch (err) {
        this.$root.$emit('notify-error', err)
      }
    }
  }
}
</script>

<style scoped lang="sass">
.wrapper
  position: relative
</style>
