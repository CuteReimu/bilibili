import os
import sys


def snake_to_camel(snake_str):
    parts = snake_str.split('_')
    return ''.join(word.title() for word in parts)


head_mapping = {
	'参数名': 'name',
	'字段': 'name',
	'类型': 'type',
	'内容': 'content',
	'必要性': 'notnull',
	'备注': 'comment',
}


if __name__ == '__main__':
	while True:
		head = []
		lines = []
		# 从控制台读取输入
		print('请输入Markdown表格，在最后一行之后输入ok表示结束：（退出请输入exit）')
		while True:
			line = input()
			if not line or line.strip() == 'exit':
				sys.exit(0)
			if line.strip() == 'ok':
				break
			if len(head) == 0:
				head = line.split('|')
				for i in range(len(head)):
					v = head[i].strip()
					head[i] = head_mapping.get(v, '')
			else:
				l = line.split('|')
				if len(l) != len(head):
					print('列数不一致，无法解析')
					sys.exit(0)
				m = {}
				for i in range(len(l)):
					v = l[i].strip()
					if v.startswith('---'):
						m = None
						break
					m[head[i]] = v
				if not m:
					continue
				lines.append(m)
		# 生成结构体
		print('type T struct {')
		for m in lines:
			name = snake_to_camel(m['name'])
			if m['type'] == 'num':
				m['type'] = 'int'
			elif m['type'] == 'str':
				m['type'] = 'string'
			elif m['type'] == 'bool':
				pass
			elif m['type'] == 'array':
				m['type'] = '[]' + name
			elif m['type'] == 'obj':
				m['type'] = name
			if m.get('notnull', '必要') == '必要':
				m['notnull'] = ''
			else:
				m['notnull'] = ',omitempty'
			content = m.get('content', '')
			comment = m.get('comment', '')
			if content or comment:
				sep = '。' if content and comment else ''
				comment = ' // %s%s%s' % (content, sep, comment)
				comment = comment.replace('<br>', '。').replace('<br/>', '。').replace('<br />', '。')
			print('\t%s %s `json:"%s%s"`%s' % (name, m['type'], m['name'], m['notnull'], comment))
		print('}')

