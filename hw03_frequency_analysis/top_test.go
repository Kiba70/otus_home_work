package hw03frequencyanalysis

import (
	"testing"

	"github.com/stretchr/testify/require" //nolint: depguard
)

// Change to true if needed.
var taskWithAsteriskIsCompleted = true

var (
	text = `Как видите, он  спускается  по  лестнице  вслед  за  своим
	другом   Кристофером   Робином,   головой   вниз,  пересчитывая
	ступеньки собственным затылком:  бум-бум-бум.  Другого  способа
	сходить  с  лестницы  он  пока  не  знает.  Иногда ему, правда,
		кажется, что можно бы найти какой-то другой способ, если бы  он
	только   мог   на  минутку  перестать  бумкать  и  как  следует
	сосредоточиться. Но увы - сосредоточиться-то ему и некогда.
		Как бы то ни было, вот он уже спустился  и  готов  с  вами
	познакомиться.
	- Винни-Пух. Очень приятно!
		Вас,  вероятно,  удивляет, почему его так странно зовут, а
	если вы знаете английский, то вы удивитесь еще больше.
		Это необыкновенное имя подарил ему Кристофер  Робин.  Надо
	вам  сказать,  что  когда-то Кристофер Робин был знаком с одним
	лебедем на пруду, которого он звал Пухом. Для лебедя  это  было
	очень   подходящее  имя,  потому  что  если  ты  зовешь  лебедя
	громко: "Пу-ух! Пу-ух!"- а он  не  откликается,  то  ты  всегда
	можешь  сделать вид, что ты просто понарошку стрелял; а если ты
	звал его тихо, то все подумают, что ты  просто  подул  себе  на
	нос.  Лебедь  потом  куда-то делся, а имя осталось, и Кристофер
	Робин решил отдать его своему медвежонку, чтобы оно не  пропало
	зря.
		А  Винни - так звали самую лучшую, самую добрую медведицу
	в  зоологическом  саду,  которую  очень-очень  любил  Кристофер
	Робин.  А  она  очень-очень  любила  его. Ее ли назвали Винни в
	честь Пуха, или Пуха назвали в ее честь - теперь уже никто  не
	знает,  даже папа Кристофера Робина. Когда-то он знал, а теперь
	забыл.
		Словом, теперь мишку зовут Винни-Пух, и вы знаете почему.
		Иногда Винни-Пух любит вечерком во что-нибудь поиграть,  а
	иногда,  особенно  когда  папа  дома,  он больше любит тихонько
	посидеть у огня и послушать какую-нибудь интересную сказку.
		В этот вечер...`

	//nolint: lll
	text2 = `Отвечать что-либо одобрительное вряд ли стоило.
	Забирать двух экспрессивных мартышек к себе на работу – тем более.
	Я не готов с таким справляться, и, в конце концов, мои родители рожали их для себя,
	а не ради того, чтоб повесить на мою шею. Решение уехать из Екатеринбурга подальше
	на юга выглядело теперь единственно правильным! К слову, никакой папки с заданием на моём столе не было.
	Возможно, потому, что в этом деле главной назначалась Гребнева, а я лишь вспомогательный резерв.
	Значит, она поставит меня в известность по всем вопросам, когда надо будет. А сейчас – завтрак!
	За накрытым столом в саду сидел я один. Как последний дурак. Герман выбежал лишь на минуту, цапнул
	полбуханки хлеба и убежал обратно. В его глазах горела жажда познания, и никакие высшие блюда любой
	кухни мира от повара с пятью звёздами Мишлен не могли бы оторвать его от золотого коня царя Митридата.
	Тут он сразу пас. Светлана даже не появилась. Но учитывая, что она ест как птичка колибри или бабочка,
	её присутствие на общекорпоративных пирах ограничивалось лишь украшением данных мероприятий.
	Она редко пила больше одного бокала, закусывала сыром или фруктами, практически не ела мяса,
	хлеба потребляла по минимуму, только белый и, как правило, с капелькой мёда. Не знаю, что это за диета,
	но фигуру она помогала сохранять идеально!       
	Может, стоит запатентовать?           
	– Саня, а ты чо тут один? – вдруг раздалось из-под стола, хотя минуту назад я сам проверял: никого там не было.
	– Где все? Куда их смыло морским прибоем из кабинета директора?
	– Вылезай! – я протянул Денисычу руку.
	– Спасибки, бро, – наш специалист по древним языкам показался на свет божий, помятый, дрожащий, сонный,
 	благоухающий винными ароматами, но, как всегда, абсолютно довольный собой.
	– Чо там было-то, куда на этот раз вострить лыжи? Я не хотел, но пришлось ответить. Может быть, это странно, но вот кому-кому,
	а именно Денисычу соврать практически невозможно! Он, конечно, махровый выпивоха, но тем не менее (или благодаря этому, хрен знает!)
	крутейший из всех известных психологов, разбирающийся в самых тонких оттенках человеческих настроений, понятий и чувств.`
)

func TestTop10(t *testing.T) {
	t.Run("no words in empty string", func(t *testing.T) {
		t.Parallel()
		require.Len(t, Top10(""), 0)
	})

	t.Run("positive test", func(t *testing.T) {
		t.Parallel()
		if taskWithAsteriskIsCompleted {
			expected := []string{
				"а",         // 8
				"он",        // 8
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"в",         // 4
				"его",       // 4
				"если",      // 4
				"кристофер", // 4
				"не",        // 4
			}
			require.Equal(t, expected, Top10(text))
		} else {
			expected := []string{
				"он",        // 8
				"а",         // 6
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"-",         // 4
				"Кристофер", // 4
				"если",      // 4
				"не",        // 4
				"то",        // 4
			}
			require.Equal(t, expected, Top10(text))
		}
	})

	t.Run("positive test 2", func(t *testing.T) {
		t.Parallel()
		if taskWithAsteriskIsCompleted {
			expected := []string{
				"не",  // 10
				"на",  // 8
				"в",   // 6
				"но",  // 6
				"я",   // 6
				"а",   // 5
				"и",   // 5
				"как", // 4
				"она", // 4
				"с",   // 4
			}
			require.Equal(t, expected, Top10(text2))
		}
	})
}
