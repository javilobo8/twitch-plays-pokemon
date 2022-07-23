const tmi = require('tmi.js');
const robot = require('robotjs');

const config = require('./config.json');

const bot = tmi.Client({
  connection: {
    cluster: "aws",
    reconnect: true
  },
  channels: [config.channel],
});

let lastMessageDate = new Date(0);

const ACTION_KEYS_BINDINGS = {
  BUTTON_UP: "up",
  BUTTON_DOWN: "down",
  BUTTON_LEFT: "left",
  BUTTON_RIGHT: "right",
  BUTTON_START: "enter",
  BUTTON_SELECT: "backspace",
  BUTTON_A: "Z",
  BUTTON_B: "X",
  BUTTON_L: "A",
  BUTTON_R: "S",
};

function handleKey(user, message) {
  console.log('handleKey', user, message);
  const actionItem = config.keys.find((i) => i.message === message);
  if (actionItem) {
    try {
      console.log(`${user}: ${message} Pressed ${actionItem.action} => ${ACTION_KEYS_BINDINGS[actionItem.action]}`)
      robot.keyTap(ACTION_KEYS_BINDINGS[actionItem.action]);
    } catch (error) {
      console.error(error);
    }
  }
}

bot.on('message', (channel, userstate, message, self) => {
  if (self) {
    return;
  }

  if (lastMessageDate.getTime() + config.interval * 1000 < new Date().getTime()) {
    lastMessageDate = new Date();
    handleKey(userstate['display-name'], message);
  }
});

bot.on('connected', () => {
  console.log('Bot connected!');
});

bot.connect();