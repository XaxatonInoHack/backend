package parser

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseReview(t *testing.T) {
	tests := []struct {
		name       string
		text       string
		wantScores map[string]float64
		wantDesc   map[string]string
	}{
		{
			name: "base",
			text: "\"Based on the reviews," +
				" I'll evaluate the employee on the given criteria." +
				" Here are my assessments:\\n\\n**Professionalism: " +
				"5/5**\\nThe employee consistently demonstrates a high " +
				"level of professionalism, as mentioned in Review 1," +
				" where they are described as having a \\\"высокий " +
				"уровень профессионализма\\\" (high level of" +
				" professionalism). Review 5 also highlights their" +
				" exceptional programming skills, calling them the " +
				"\\\"лучший программист эвер\\\" (best programmer ever)." +
				" The employee's ability to provide logical explanations" +
				" for their decisions, as mentioned in Review 1," +
				" further supports their professionalism. Their expertise is " +
				"evident in their work, and they maintain a professional demeanor." +
				" Overall, the employee's professionalism is outstanding." +
				"\\n\\n**Teamwork: 4.5/5**\\nReview 4 highlights the employee's " +
				"\\\"командную работу\\\" (teamwork), indicating that they are" +
				" able to collaborate effectively with others. Additionally," +
				" Review 1 mentions that the employee is always willing to help and " +
				"provide answers to questions, which suggests that they are a team player." +
				" However, Review 3 expresses a desire for more proactive approach," +
				" which might indicate some room for improvement in this area. Nevertheless," +
				" the overall sentiment suggests that the employee is a strong team player. " +
				"The only reason I'm not giving a 5 is that there is a faint criticism in Review " +
				"3.\\n\\n**Communication: 5**\\nThe employee's communication skills are" +
				" consistently praised across the reviews. Review 1 highlights their ability" +
				" to provide clear and logical explanations, while Review 4 mentions their" +
				" \\\"открытость к диалогу\\\" (openness to dialogue). The employee's " +
				"responsiveness and willingness to help, as mentioned in Review 1, also " +
				"demonstrate their strong communication skills. Additionally, Review 5" +
				" expresses appreciation for the employee's excellent communication style. " +
				"Overall, the employee's communication skills are exceptional." +
				"\\n\\n**Initiative: 3/5**\\nWhile Review 4 mentions the employee's " +
				"\\\"нацеленность на результат\\\" (focus on results), Review 3\"",
			wantScores: map[string]float64{"Communication": 5, "Professionalism": 5, "Teamwork": 4.5, "Initiative": 3},
		},
		{
			name:       "all",
			text:       "Based on the reviews, I'll evaluate the employees on the given criteria. Please note that the reviews are not labeled with specific employee names, so I'll provide a general evaluation for each employee ID.\\n\\n**Employee ID - 59595**\\n1. Professionalism: 5\\nThe employee demonstrates deep expertise in technical and organizational aspects of the company and can find solutions to complex situations. Their knowledge and experience are highly valued by colleagues. They provide valuable feedback and suggestions, and their involvement in projects is highly appreciated. Their professionalism is evident in their ability to balance security, timelines, and comfort of work. Overall, they are considered an expert in their field.\\n\\n2. Teamwork: 5\\nThe employee is praised for their ability to work collaboratively with colleagues, finding compromises and solutions that benefit everyone. They are willing to help and provide guidance, making them a valuable team player. Their involvement in team projects is highly appreciated, and they are able to bring people together to achieve common goals. They demonstrate a strong sense of responsibility and accountability. Their team-oriented approach is evident in their ability to find mutually beneficial solutions.\\n\\n3. Communication: 5\\nThe employee is commended for their excellent communication skills, being able to explain complex issues in a clear and concise manner. They are approachable, open, and willing to listen to others. Their communication style is constructive, and they are able to find common ground with colleagues. They are able to articulate their thoughts and opinions effectively, making them an effective communicator. Their ability to communicate technical information to non-technical colleagues is particularly valued.\\n\\n4. Initiative: 5\\nThe employee is praised for their proactive approach, taking the initiative to find solutions to problems and improve processes. They are willing to take on new challenges and are not afraid to think outside the box. Their creativity and resourcefulness are highly valued by colleagues. They demonstrate a strong sense of ownership and accountability, taking charge of projects and seeing them through to completion. Their willingness to",
			wantDesc:   map[string]string{},
			wantScores: map[string]float64{},
		},
		{
			name:       "short",
			text:       "Wow**Proffesionalism:5**gdsfd. ahhsdfh**Ink: 4.5**, gaa",
			wantScores: map[string]float64{"Proffesionalism": 5, "Ink": 4.5},
			wantDesc:   map[string]string{"Proffesionalism": "gdsfd. ahhsdfh", "Ink": ", gaa"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scores, desc := ParseReview(tt.text)
			if tt.wantScores != nil {
				require.Equal(t, tt.wantScores, scores)
			}
			if tt.wantDesc != nil {
				require.Equal(t, tt.wantDesc, desc)
			}

		})
	}
}
