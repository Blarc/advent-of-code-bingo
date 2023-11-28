import {BingoBoardDto} from './bingoBoardDto';
import {BingoCardDto} from './bingoCardDto';

export class UserDto {
    'avatar_url': string;
    'name': string;
    'github_url': string;
    'reddit_url': string;
    'bingo_cards': BingoCardDto[];
    'bingo_boards': BingoBoardDto[];
}
