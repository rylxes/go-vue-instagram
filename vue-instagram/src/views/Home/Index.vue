<template>
  <v-layout>
    <v-card contextual-style="dark">
      <span slot="header">
        {{ $t('general.welcome') }}
      </span>
      <div slot="body">
        <p>
          Calculate Instagram average engagement rate.
        </p>

        <div slot="body">
          <form @submit.prevent="calculator(user)">
            <div class="form-group">
              <div class="input-group">
                <div class="input-group-prepend">
                  <span class="input-group-text">
                    <i class="fa fa-envelope fa-fw" />
                  </span>
                </div>
                <input
                  v-model="user.username"
                  type="text"
                  required
                  placeholder="Instagram Username"
                  class="form-control"
                >
              </div>
            </div>
            <div class="form-group">
              <button class="btn btn-outline-primary">
                Calculate
              </button>
            </div>
          </form>
        </div>
        <span v-if="average">
          {{ average }} %
        </span>

        <span v-if="error" style="background-color: red">
          {{ error }}
        </span>
      </div>
      <div slot="footer">
        Made by rylxes
      </div>
    </v-card>
  </v-layout>
</template>

<script>
/* ============
 * Home Index Page
 * ============
 *
 * The home index page.
 */
import axios from 'axios';
import VLayout from '@/layouts/Default.vue';
import VCard from '@/components/Card.vue';

export default {
  /**
   * The name of the page.
   */
  name: 'HomeIndex',

  /**
   * The components that the page can use.
   */
  components: {
    VLayout,
    VCard,
  },

  /**
   * The data that can be used by the page.
   *
   * @returns {Object} The view-model data.
   */
  data() {
    return {
      average: null,
      error: null,
      user: {
        username: null,
      },
    };
  },

  /**
   * The methods the page can use.
   */
  methods: {
    /**
     * Will log the user in.
     *
     * @param {Object} user The user to be logged in.
     */
    calculator(user) {
      console.log(user);
      this.average = null;
      this.error = null;
      const baseURI = `http://127.0.0.1:8000/calculate/${user.username}`;
      axios.get(baseURI)
        .then((result) => {
          console.log(result.data);
          this.average = result.data;
        }).catch((error) => {
          this.error = error;
          console.log(error);
        });
    },
  },
};
</script>
