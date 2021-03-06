package gameprocessor

const (
	responseGamePaused                  = "Участие в игре приостановлено. Ничего не сможешь делать, пока не напишешь " + commandResume + "."
	responseGameResumed                 = "Участие в игре возобновлено. Продолжай как ни в чем не бывало."
	responseAskName                     = "Привет! Чтобы стать участником Свитчеров, напиши в ответ свое имя. Важно, чтобы другие участники могли тебя узнать, так что не пиши ерунду."
	responseNiceToMeet                  = "Приятно познакомиться, %s. Теперь жди инструкции. Они могут приходить в любой момент, так что держи телефон включенным! Чтобы приостановить участие в игре, напиши " + commandPause + "."
	responseSetName                     = "Напиши свое имя. Важно, чтобы другие участники могли тебя узнать, так что не пиши ерунду."
	responseLeaders                     = "Вот кто заработал больше всего очков в Свитчерах:\n\n%s\nА у тебя к этому времени накопилось всего %d."
	responsePlayerGathered              = "Ждем отстающих еще немного и начинаем."
	responseGatheringInstructions       = "Соберитесь в указанном месте. Как только соберетесь, все должны написать \"" + commandGathered + "\". Постарайтесь дождаться всю команду, прежде чем писать \"" + commandGathered + "\", но не забывайте о дедлайне."
	responseGatherNotAnswer             = "Все уже собрались, можно больше не писать \"" + commandGathered + "\". Отвечай на задание."
	responsePlayerAnswered              = "Ответ принят."
	responseWaitForModeration           = "Ответ на модерации, жди решение."
	responseSomethingWrong              = "Что-то пошло не так. Попробуй еще раз."
	responseDefault                     = "Жди инструкции или напиши какую-нибудь команду. Я понимаю:\n\n" + commandSetName + " — изменить имя\n" + commandPause + " — приостановить участие в игре\n" + commandLeaders + " — посмотреть лидеров"
	responsePlayerFailedToGather        = "Нужно было вовремя написать \"" + commandGathered + "\", а у тебя не получилось. Жди теперь следующий раунд."
	responseTeamFailedToGather          = "Время вышло :( Этот раунд вы проиграли, потому что не собрали команду вовремя. Но в следующий раз повезет! Ждите следующий раунд."
	responseTeamWon                     = "Вы победили и получаете кучу очков! Ждите следующий раунд."
	responseTeamLost                    = "Вы проиграли, потому что ответили неправильно. Теперь вы не получите кучу очков :("
	responseModerationRequired          = "Ваш ответ направлен на модерацию, ждите решения."
	responseGatheringTaskSuffix         = " Как только соберетесь, напишите \"" + commandGathered + "\" в чат. Постарайтесь дождаться всю команду, прежде чем писать \"" + commandGathered + "\", но не забывайте о дедлайне."
	responseActualTaskSuffix            = " Ответ может прислать любой участник, шансов передумать (почти) нет. Если несколько ответов пришли приблизительно одновременно, засчитывается последний."
	responseTeamFailedToAnswer          = "Время вышло :( Этот раунд вы проиграли, потому что не ответили на задачу вовремя. Но в следующий раз повезет! Ждите следующий раунд."
	responseTrumpActiveRoundFinished    = "Активный раунд завершился."
	responseTrumpTeamGotQuorum          = "Команда %d набрала кворум, даем задачу."
	responseTrumpTeamFailedToGather     = "У команды %d закончилось время на сборы, они проиграли."
	responseTrumpTeamWon                = "Команда %d выиграла, дав правильный ответ."
	responseTrumpTeamLost               = "Команда %d проиграла, дав неправильный ответ."
	responseTrumpModerationRequired     = "Команда %d дала ответ, требуется модерация."
	responseTrumpTeamFailedToAnswer     = "У команды %d закончилось время на ответ, они проиграли."
	responseTrumpNothingToModerate      = "Сейчас нет активного раунда, модерировать нечего."
	responseTrumpSomethingWrong         = "Произошла ошибка: %s"
	responseTrumpAlreadyModerated       = "Эту команду кто-то уже отмодерировал."
	responseTrumpRoundStarted           = "Начался новый раунд."
	responseTrumpTaskPrefix             = "Задание: "
	responseTrumpModerationInstructions = "Напиши \"" + commandYes + "\", если ответ правильный. Напиши \"" + commandNo + "\", если ответ неправильный. Напиши что угодно другое, чтобы бросить это занятие и пойти строить стену на границе с Мексикой."
	responseTrumpResigned               = "Отставка принята."
	responseTrumpDefault                = "Добро пожаловать, господин президент. Издайте какой-нибудь указ:\n\n" + commandNewRound + " — запустить новый раунд\n" + commandModerate + " — модерировать что-нибудь\n" + commandLeaders + " — посмотреть очки игроков\n" + commandResign + " — подать в отставку"

	teamStateGathering  = "gathering"
	teamStatePlaying    = "playing"
	teamStateModeration = "moderation"
	teamStateWon        = "won"
	teamStateLost       = "lost"

	commandNewRound = "/newround"
	commandModerate = "/moderate"
	commandResign   = "/resign"
	commandSetName  = "/setname"
	commandPause    = "/pause"
	commandResume   = "/resume"
	commandLeaders  = "/leaders"
	commandGathered = "тут"
	commandYes      = "да"
	commandNo       = "нет"
)
