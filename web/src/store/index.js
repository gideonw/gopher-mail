import Vue from 'vue'
import Vuex, { Store } from 'vuex'
import axios from "axios";
import _ from 'lodash';

import {
  LIST_EMAILS,
  LOGIN,
  LOGOUT,
  GET_EMAIL,
} from './actions'

import {
  M_LOGIN,
  M_LOGOUT,
  M_STORE_EMAILS,
} from './mutations'

Vue.use(Vuex)

const base_domain = "https://gps.gideonw.xyz"

export default new Store({
  state: {
    auth: {
      valid: false,
      username: '',
      login: {
        state: 'NONE'
      }
    },
    emails: {}
  },
  mutations: {
    [M_LOGIN](state) {
      state.auth.valid = true;
    },
    [M_LOGOUT](state) {
      state.auth.valid = false;
    },
    [M_STORE_EMAILS](state, emails) {
      if (emails.messageID) {
        // single email
        state.emails[emails.messageID] = emails.email;
      } else {
        // email list
        _.forEach(emails.email, function (value) {
          state.emails[value] = {};
        });

      }
    }
  },
  actions: {
    [LOGIN](context) {
      context.commit(M_LOGIN);
    },
    [LOGOUT](context) {
      context.commit(M_LOGOUT);
    },
    [LIST_EMAILS](context) {
      axios.get(`${base_domain}/api/gideon/emails`).then(resp => {
        console.log(resp);
        context.dispatch(M_STORE_EMAILS, resp.data.emails);
      });
    },
    [GET_EMAIL](context, emailID) {
      axios.get(`${base_domain}/api/gideon/email/${emailID}`).then(resp => {
        console.log(resp);
        context.dispatch(M_STORE_EMAILS, resp.data);
      });
    }
  },
  modules: {
  }
})
