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
  M_STORE_EMAIL,
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
    [M_STORE_EMAIL](state, emails) {
      console.log(emails);
      if (emails.messageID) {
        // single email
        state.emails[emails.messageID] = emails.email;
      } else {
        // email list
        _.forEach(emails, function (value) {
          state.emails[value] = {};
        });
      }
    }
  },
  actions: {
    [LOGIN](context) {
      context.commit(M_LOGIN);
      console.log("dispatch");
      context.dispatch(LIST_EMAILS, null, { root: true });
    },
    [LOGOUT](context) {
      context.commit(M_LOGOUT);
    },
    [LIST_EMAILS](context) {
      axios.get(`${base_domain}/api/gideon/emails`).then(resp => {
        context.commit(M_STORE_EMAIL, resp.data.emails);
        _.forEach(resp.data.emails, function (value) {
          context.dispatch(GET_EMAIL, value, { root: true });
        });
      });
    },
    [GET_EMAIL](context, emailID) {
      axios.get(`${base_domain}/api/gideon/email/${emailID}`).then(resp => {
        context.commit(M_STORE_EMAIL, resp.data);
      });
    }
  },
  modules: {
  }
})
