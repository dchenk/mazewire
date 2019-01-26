import Quill from "quill/core/quill"; // the core Quill file, without any modules or formats included in it

import Block, {BlockEmbed} from "quill/blots/block";
import Break from "quill/blots/break";
import Container from "quill/blots/container";
import Cursor from "quill/blots/cursor";
import Embed from "quill/blots/embed";
import Inline from "quill/blots/inline";
import Scroll from "quill/blots/scroll";
import TextBlot from "quill/blots/text";

import Clipboard from "quill/modules/clipboard";
import History from "quill/modules/history";
import Keyboard from "quill/modules/keyboard";
import Toolbar from "quill/modules/toolbar";

import Bold from "quill/formats/bold";
import Italic from "quill/formats/italic";
import Header from "quill/formats/header";
import Underline from "quill/formats/underline";
import Link from "quill/formats/link";
// import CodeBlock, {Code as InlineCode} from "quill/formats/code";
import CodeBlock from "quill/formats/code";
import List, {ListItem} from "quill/formats/list";
import Script from "quill/formats/script";
import {ColorStyle} from "quill/formats/color";
import {AlignClass} from "quill/formats/align";
import {IndentClass as Indent} from "quill/formats/indent";

import ColorPicker from "quill/ui/color-picker";

import Parchment from "quill/node_modules/parchment"; // get from here to not double-import

import Snow from "quill/themes/snow";
import "quill/dist/quill.snow.css";

Quill.register({

	"blots/block"       	: Block,
	"blots/block/embed" 	: BlockEmbed,
	"blots/break"       	: Break,
	"blots/container"   	: Container,
	"blots/cursor"      	: Cursor,
	"blots/embed"       	: Embed,
	"blots/inline"      	: Inline,
	"blots/scroll"      	: Scroll,
	"blots/text"			: TextBlot,

	"modules/clipboard"		: Clipboard,
	"modules/history"  		: History,
	"modules/keyboard" 		: Keyboard,
	"modules/toolbar"		: Toolbar,

	"themes/snow"			: Snow,

	"formats/bold"			: Bold,
	"formats/italic"		: Italic,
	"formats/header"		: Header,
	"formats/underline"		: Underline,
	"formats/list"			: List,
	"formats/list/item"		: ListItem,
	"formats/link"			: Link,
	"formats/script"		: Script,
	"formats/code-block"	: CodeBlock,
	"formats/color"			: ColorStyle,
	"formats/align"			: AlignClass,
	"formats/indent"		: Indent,

	"ui/color-picker"		: ColorPicker

});

Parchment.register(Block, Break, Cursor, Inline, Scroll, TextBlot, ListItem, CodeBlock);

export default Quill;
