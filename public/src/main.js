
let SERVER = "http://52.27.58.29:3000";
var UUID;

var chatData = {
	UUID: '',
	newChat: '',
	chats: [],
	timer: null,
	timerIsOn: false,
	timedResponseInterval: 0,
}

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
		scrollToBottom: function () {
			setTimeout(function () {  
				items = document.querySelectorAll(".chat");
				console.log(items);
		    	last = items[items.length-1];
		    	last.scrollIntoView();
    		}, 30);
		},
		pushUserChat: function () {

			chatData.chats.push({
				text: chatData.newChat,
				isBot: false,
				name: "You",
				date: new Date().toString()
			});
		},
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
		submitChat: function (e) {
			this.addNewChat(UUID);
		}, 
		addNewChat: addNewChat = function (UUIDval) {

			axios.get('/conversation/' + chatData.UUID + "?text=" + chatData.newChat)
			.then(function (response) {
				chatData.timedResponseInterval = response.data.Timed_Response_Interval;
				console.log(chatData.timedResponseInterval);
				if (chatData.newChat.length) {
					console.log("canceling timer");
					clearTimeout(chatData.timer);
					chatData.timerIsOn = false;
					console.log("pushing chat");
					vm.pushUserChat();
				}
				chatData.newChat = '';
				vm.pushBotChat(response);
				chatData.newChat = '';
				chatData.UUID = response.data.Conversation;
				if (!chatData.timerIsOn && chatData.timedResponseInterval) {
				console.log("timer started");
				chatData.timer = setTimeout(function () {
						chatData.timerIsOn = true;
						console.log("issuing timer request");
	 					axios.get('/wakeup/' + chatData.UUID)
	 					.then(function (response) {
	 						console.log("pushing timer chat");
	 						vm.pushBotChat(response);
	 						vm.scrollToBottom();
	 					})
	 						.catch(function (error) {
	 						console.log(error);
	 					});
					}, response.data.Timed_Response_Interval * 1000);
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






