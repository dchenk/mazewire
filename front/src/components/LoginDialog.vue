<template>
	<div id="login-dialog">
		<p class="err-message" v-if="errMsg !== ''">
			{{ errMsg }}
		</p>
		<div class="dyn-textfield">
			<input type="text" class="dyn-textfield-input" id="login-user" @keyup.enter="$parent.accept" v-model="user">
			<label for="login-user" class="dyn-textfield-label">Username or Email</label>
		</div>
		<div class="dyn-textfield">
			<input type="password" class="dyn-textfield-input" id="login-pass" @keyup.enter="$parent.accept" v-model="pass">
			<label for="login-pass" class="dyn-textfield-label">Password</label>
		</div>
	</div>
</template>

<script>

	import DynamicTextFields from "dynamic-textfields"

	export default {
		props: {
			method: {
				type: String,
				default: ""
			},
			endpoint: {
				type: String,
				default: ""
			},
			data: {
				type: Object,
				default: null
			},
			handleGood: {
				type: Function,
				default: () => {}
			},
			handleError: {
				type: Function,
				default: () => {}
			}
		},
		data() { return {
			user: "",
			pass: "",
			errMsg: ""
		}},
		mounted() {
			(new DynamicTextFields("login-dialog")).registerAll()
		},
		methods: {
			acceptResponse() {
				this.$req("POST", "auth",
					{
						"user": this.user,
						"pass": this.pass
					},
					_ => { // On success
						this.$hideDialog()
						this.$toasted.info("Retrying request...", {duration: 1200})
						if (this.$parent.responseCallback) {
							this.$parent.responseCallback()
						}
					},
					err => {
						this.errMsg = err
					}
				)
			}
		}
	}

</script>
