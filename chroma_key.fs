#version 330

in vec2 fragTexCoord;
in vec4 fragColor;

out vec4 finalColor;

uniform sampler2D texture0;
uniform vec3 keyColor; // Background color to remove (normalized RGB, e.g., black = 0.0, 0.0, 0.0)
uniform float threshold; // Tolerance for color matching

void main()
{
    vec4 texColor = texture(texture0, fragTexCoord) * fragColor;
    vec3 diff = abs(texColor.rgb - keyColor);
    if (max(max(diff.r, diff.g), diff.b) < threshold)
    {
        discard; // Discard pixel if it matches the key color
    }
    finalColor = texColor;
}