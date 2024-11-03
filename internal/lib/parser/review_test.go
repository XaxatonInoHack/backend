package parser

import (
	"reflect"
	"testing"
)

func TestParseScores(t *testing.T) {
	text := `Based on the reviews, I'll evaluate the employee's performance as follows:

1. Professionalism: 5
The employee is described as having high professionalism in Review 1, with traits like intelligence, politeness, and the ability to balance interests. Review 5 also praises the employee as the "лучший программист эвер" (best programmer ever), suggesting exceptional professional skills. This is consistent with the evaluation from others, who gave a high score for Initiative. Although Review 3 hints at some potential issues, it's a relatively minor criticism and doesn't seem to impugn the employee's professionalism. Overall, the reviews suggest a high level of professionalism.

2. Teamwork: 4.5
Review 4 praises the employee's "командную работу" (teamwork), indicating a strong ability to collaborate with others. Review 1 also describes the employee as someone who can be approached with any questions, suggesting openness and a willingness to help. While Review 3 expresses a minor criticism, it's not directly related to teamwork. The employee's own evaluation score for Teamwork is high, which is consistent with the reviews. Overall, the reviews suggest strong teamwork skills.

3. Communication: 4.5
Review 1 highlights the employee's excellent communication skills, praising their ability to provide clear explanations. Review 4 also appreciates the employee's "открытость к диалогу" (openness to dialogue), indicating strong communication. Although Review 3 expresses a minor criticism, it's not directly related to communication. The reviews suggest that the employee is generally clear and approachable. The employee's own evaluation score for Communication is relatively low, but the reviews suggest this might be an underestimation.

4. Initiative: 5
The employee's own evaluation score for Initiative is already high, and Review 1 suggests a proactive approach to helping others. Review 4 praises the employee's "нацеленность на результат" (focus`

	expected := map[string]float64{
		"Professionalism": 5.0,
		"Teamwork":        4.5,
		"Communication":   4.5,
		"Initiative":      5.0,
	}

	result := ParseScoreOnly(text)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Ожидается %v, получено %v", expected, result)
	}
}

func TestParseScoresWithMissingData(t *testing.T) {
	text := `1. Professionalism: 4.8
Some text without other categories.`

	expected := map[string]float64{
		"Professionalism": 4.8,
	}

	result := ParseScoreOnly(text)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Ожидается %v, получено %v", expected, result)
	}
}

func TestParseScoresWithInvalidData(t *testing.T) {
	text := `1. Professionalism: five
This text contains invalid score.`

	expected := map[string]float64{}

	result := ParseScoreOnly(text)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Ожидается %v, получено %v", expected, result)
	}
}

func TestParseCategoryTexts(t *testing.T) {
	text := `Based on the reviews, I'll evaluate the employees on the given criteria. Please note that the reviews are not labeled with specific employee names, so I'll provide a general evaluation for each employee ID.

**Employee ID - 59595**
1. Professionalism: 5
The employee demonstrates deep expertise in technical and organizational aspects of the company and can find solutions to complex situations. Their knowledge and experience are highly valued by colleagues. They provide valuable feedback and suggestions, and their involvement in projects is highly appreciated. Their professionalism is evident in their ability to balance security, timelines, and comfort of work. Overall, they are considered an expert in their field.

2. Teamwork: 5
The employee is praised for their ability to work collaboratively with colleagues, finding compromises and solutions that benefit everyone. They are willing to help and provide guidance, making them a valuable team player. Their involvement in team projects is highly appreciated, and they are able to bring people together to achieve common goals. They demonstrate a strong sense of responsibility and accountability. Their team-oriented approach is evident in their ability to find mutually beneficial solutions.

3. Communication: 5
The employee is commended for their excellent communication skills, being able to explain complex issues in a clear and concise manner. They are approachable, open, and willing to listen to others. Their communication style is constructive, and they are able to find common ground with colleagues. They are able to articulate their thoughts and opinions effectively, making them an effective communicator. Their ability to communicate technical information to non-technical colleagues is particularly valued.

4. Initiative: 5
The employee is praised for their proactive approach, taking the initiative to find solutions to problems and improve processes. They are willing to take on new challenges and are not afraid to think outside the box. Their creativity and resourcefulness are highly valued by colleagues. They demonstrate a strong sense of ownership and accountability, taking charge of projects and seeing them through to completion. Their willingness to`

	expected := map[string]string{
		"Professionalism": `The employee demonstrates deep expertise in technical and organizational aspects of the company and can find solutions to complex situations. Their knowledge and experience are highly valued by colleagues. They provide valuable feedback and suggestions, and their involvement in projects is highly appreciated. Their professionalism is evident in their ability to balance security, timelines, and comfort of work. Overall, they are considered an expert in their field.`,
		"Teamwork":        `The employee is praised for their ability to work collaboratively with colleagues, finding compromises and solutions that benefit everyone. They are willing to help and provide guidance, making them a valuable team player. Their involvement in team projects is highly appreciated, and they are able to bring people together to achieve common goals. They demonstrate a strong sense of responsibility and accountability. Their team-oriented approach is evident in their ability to find mutually beneficial solutions.`,
		"Communication":   `The employee is commended for their excellent communication skills, being able to explain complex issues in a clear and concise manner. They are approachable, open, and willing to listen to others. Their communication style is constructive, and they are able to find common ground with colleagues. They are able to articulate their thoughts and opinions effectively, making them an effective communicator. Their ability to communicate technical information to non-technical colleagues is particularly valued.`,
		"Initiative":      `The employee is praised for their proactive approach, taking the initiative to find solutions to problems and improve processes. They are willing to take on new challenges and are not afraid to think outside the box. Their creativity and resourcefulness are highly valued by colleagues. They demonstrate a strong sense of ownership and accountability, taking charge of projects and seeing them through to completion. Their willingness to`,
	}

	result := ParseCategoryTexts(text)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Ожидается %v, получено %v", expected, result)
	}
}

func TestParseCategoryTextsWithMultipleEntries(t *testing.T) {
	text := `
1. Professionalism: 4.5
First professionalism text.

2. Teamwork: 4.0
First teamwork text.

1. Professionalism: 5
Second professionalism text.

2. Communication: 4.8
Communication text.

5. Overall Perfomance: 4.8
Text text text
`

	expected := map[string]string{
		"Professionalism":    "Second professionalism text.",
		"Teamwork":           "First teamwork text.",
		"Communication":      "Communication text.",
		"Overall Perfomance": "Text text text",
	}

	result := ParseCategoryTexts(text)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Ожидается %v, получено %v", expected, result)
	}
}

func TestParseCategoryTextsWithNoCategories(t *testing.T) {
	text := `This text does not contain any categories that match the pattern.`

	expected := map[string]string{}

	result := ParseCategoryTexts(text)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Ожидается пустая карта, получено %v", result)
	}
}

func TestParseCategoryTextsWithIncompleteData(t *testing.T) {
	text := `
1. Professionalism:
No score provided.

2. Teamwork: 4.5
Teamwork text.
`

	expected := map[string]string{
		"Teamwork": "Teamwork text.",
	}

	result := ParseCategoryTexts(text)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Ожидается %v, получено %v", expected, result)
	}
}
