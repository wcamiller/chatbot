
/* Vue.js data object */

var chatData = {
	UUID: '',
	newChat: '',
	chats: [],
	timer: null,
	timerIsOn: false,
	timedResponseInterval: 0,
}

/*  HTML chat template */

Vue.component('chat', {
	template: 		  '<li><div><span>' + 
					  '<small><b>{{ name }}</b> - {{ date }}</small></span><br><br>' + 
					  '<p>{{ text }}</p></div></li>',
	props: ['text', 'name', 'date']
});

var vm = new Vue({
	el: '#el',
	data: chatData,
	methods: {

		/* helper method that issues a POST request to the /wakeup endpoint (which simply sends an empty JSON to 
		the /coversation/<UUID> endpoint of the PullString Web API) based on the timedResponseInterval of the 
		previous conversation response (if any) */

		setInputTimer: function () {
			chatData.timer = setTimeout(function () {
				chatData.timerIsOn = true;
	 			axios.get('/wakeup/' + chatData.UUID)
	 			.then(function (response) {
	 				vm.pushBotChat(response);
	 				vm.scrollToBottom();
	 			})
	 			.catch(function (error) {
	 				console.log(error);
	 			});
			}, chatData.timedResponseInterval * 1000);
		},

		/* method to advance chat window upon receipt of new chat message */

		scrollToBottom: function () {
			setTimeout(function () {  
				items = document.querySelectorAll(".chat");
		    	last = items[items.length-1];
		    	last.scrollIntoView();
    		}, 30);
		},

		/* pushes user-inputted chat msg to Vue data object along with date and name information as well as boolean
		to trigger dynamic CSS class assignments */

		pushUserChat: function () {

			chatData.chats.push({
				text: chatData.newChat,
				isBot: false,
				name: "You",
				date: new Date().toString()
			});
		},

		/* pushes chat msg received in JSON response of PullString Web API along with date, name and CSS class info */

		pushBotChat: function (response) {
			var arr = response.data.Outputs;
			for (var i in arr) {
				chatData.chats.push({
					text: arr[i].Text,
					isBot: true,
					name: "PullString Bot",
					date: new Date().toString()
				});
			}
		},

		/* function which is bound to the input element and is triggered upon text submission events */

		submitChat: function (e) {
			this.addNewChat(chatData.UUID);
		}, 

		/* main logic for making calls to golang server a receiving JSON data */
		addNewChat: addNewChat = function (UUIDval) {

			axios.get('/conversation/' + chatData.UUID + "?text=" + chatData.newChat)
			.then(function (response) {
				chatData.timedResponseInterval = response.data.Timed_Response_Interval;

				if (chatData.newChat.length) {  /* suppresses blank chat msg when chatData.newChat is initially empty */ 
					clearTimeout(chatData.timer);
					chatData.timerIsOn = false; /* turn off timer for /wakeup call, b/c user has entered txt */
					vm.pushUserChat();
				}
				chatData.newChat = '';
				vm.pushBotChat(response);
				chatData.newChat = '';
				chatData.UUID = response.data.Conversation;
				/* if the timer is not running and there was a timedResponseInterval associated with the previous
				chatbot response, start the timer */
				if (!chatData.timerIsOn && chatData.timedResponseInterval) {
					vm.setInputTimer();
				}
				vm.scrollToBottom();

			})
			.catch(function (error) {
				console.log(error);
			});
		}
	}
});

addNewChat("");






